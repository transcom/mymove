import React, { useState } from 'react';
import { arrayOf, bool, func, number, shape, string, oneOf } from 'prop-types';
import { Field, Formik } from 'formik';
import { generatePath } from 'react-router';
import { queryCache, useMutation } from 'react-query';
import { Alert, Button, Checkbox, Fieldset, FormGroup, Radio } from '@trussworks/react-uswds';

import getShipmentOptions from '../../Customer/MtoShipmentForm/getShipmentOptions';

import styles from './ShipmentForm.module.scss';

import { MTO_SHIPMENTS } from 'constants/queryKeys';
import { SCRequestShipmentCancellationModal } from 'components/Office/ServicesCounseling/SCRequestShipmentCancellationModal/SCRequestShipmentCancellationModal';
import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import ShipmentAccountingCodes from 'components/Office/ShipmentAccountingCodes/ShipmentAccountingCodes';
import ShipmentWeightInput from 'components/Office/ShipmentWeightInput/ShipmentWeightInput';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import StorageFacilityInfo from 'components/Office/StorageFacilityInfo/StorageFacilityInfo';
import StorageFacilityAddress from 'components/Office/StorageFacilityAddress/StorageFacilityAddress';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { servicesCounselingRoutes, tooRoutes } from 'constants/routes';
import { dropdownInputOptions } from 'shared/formatters';
import { formatWeight } from 'utils/formatters';
import { shipmentDestinationTypes } from 'constants/shipments';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { HhgShipmentShape, MtoShipmentShape } from 'types/customerShapes';
import { formatMtoShipmentForAPI, formatMtoShipmentForDisplay } from 'utils/formatMtoShipment';
import { MatchShape } from 'types/officeShapes';
import { AccountingCodesShape } from 'types/accountingCodes';
import { validateDate } from 'utils/validation';
import { deleteShipment } from 'services/ghcApi';
import { officeRoles, roleTypes } from 'constants/userRoles';
import ShipmentFormRemarks from 'components/Office/ShipmentFormRemarks/ShipmentFormRemarks';
import ShipmentVendor from 'components/Office/ShipmentVendor/ShipmentVendor';

const ShipmentForm = ({
  match,
  history,
  newDutyStationAddress,
  selectedMoveType,
  isCreatePage,
  isForServicesCounseling,
  mtoShipment,
  submitHandler,
  mtoShipments,
  serviceMember,
  currentResidence,
  moveTaskOrderID,
  TACs,
  SACs,
  userRole,
  displayDestinationType,
}) => {
  const [errorMessage, setErrorMessage] = useState(null);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);

  const shipments = mtoShipments;

  const [mutateMTOShipmentStatus] = useMutation(deleteShipment, {
    onSuccess: (_, variables) => {
      const updatedMTOShipment = mtoShipment;
      // Update mtoShipments with our updated status and set query data to match
      shipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryCache.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestCancellation button to update immediately
      queryCache.invalidateQueries([MTO_SHIPMENTS, variables.moveTaskOrderID]);

      history.goBack();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setErrorMessage(errorMsg);
    },
  });

  const getShipmentNumber = () => {
    // TODO - this is not supported by IE11, shipment number should be calculable from Redux anyways
    // we should fix this also b/c it doesn't display correctly in storybook
    const { search } = window.location;
    const params = new URLSearchParams(search);
    const shipmentNumber = params.get('shipmentNumber');
    return shipmentNumber;
  };

  const handleDeleteShipment = (shipmentID) => {
    mutateMTOShipmentStatus({
      shipmentID,
    });
  };

  const handleShowCancellationModal = () => {
    setIsCancelModalVisible(true);
  };

  const shipmentType = mtoShipment.shipmentType || selectedMoveType;
  const { showDeliveryFields, showPickupFields, schema } = getShipmentOptions(shipmentType, userRole);

  const isHHG = shipmentType === SHIPMENT_OPTIONS.HHG;
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
  const showAccountingCodes = isNTS || isNTSR;

  const isTOO = userRole === roleTypes.TOO;
  const isServiceCounselor = userRole === roleTypes.SERVICES_COUNSELOR;

  const shipmentDestinationAddressOptions = dropdownInputOptions(shipmentDestinationTypes);

  const shipmentNumber = isHHG ? getShipmentNumber() : null;
  const initialValues = formatMtoShipmentForDisplay(
    isCreatePage
      ? { userRole, shipmentType }
      : { userRole, shipmentType, agents: mtoShipment.mtoAgents, ...mtoShipment },
  );
  const optionalLabel = <span className={formStyles.optional}>Optional</span>;
  const { moveCode } = match.params;

  const moveDetailsRoute = isTOO ? tooRoutes.MOVE_VIEW_PATH : servicesCounselingRoutes.MOVE_VIEW_PATH;
  const moveDetailsPath = generatePath(moveDetailsRoute, { moveCode });

  const editOrdersRoute = isTOO ? tooRoutes.ORDERS_EDIT_PATH : servicesCounselingRoutes.ORDERS_EDIT_PATH;
  const editOrdersPath = generatePath(editOrdersRoute, { moveCode });

  const submitMTOShipment = ({
    shipmentOption,
    pickup,
    hasDeliveryAddress,
    delivery,
    customerRemarks,
    counselorRemarks,
    ntsRecordedWeight,
    tacType,
    sacType,
    serviceOrderNumber,
    storageFacility,
    usesExternalVendor,
    destinationType,
  }) => {
    const deliveryDetails = delivery;
    if (hasDeliveryAddress === 'no' && shipmentType !== SHIPMENT_OPTIONS.NTSR) {
      delete deliveryDetails.address;
    }

    let nullableTacType = tacType;
    let nullableSacType = sacType;
    if (showAccountingCodes && !isCreatePage) {
      nullableTacType = typeof tacType === 'undefined' ? '' : tacType;
      nullableSacType = typeof sacType === 'undefined' ? '' : sacType;
    }

    const pendingMtoShipment = formatMtoShipmentForAPI({
      shipmentType: shipmentOption || selectedMoveType,
      moveCode,
      customerRemarks,
      counselorRemarks,
      pickup,
      delivery: deliveryDetails,
      ntsRecordedWeight,
      tacType: nullableTacType,
      sacType: nullableSacType,
      serviceOrderNumber,
      storageFacility,
      usesExternalVendor,
      destinationType,
    });

    const updateMTOShipmentPayload = {
      moveTaskOrderID,
      shipmentID: mtoShipment.id,
      ifMatchETag: mtoShipment.eTag,
      normalize: false,
      body: pendingMtoShipment,
    };

    if (isCreatePage) {
      const body = { ...pendingMtoShipment, moveTaskOrderID };
      submitHandler({ body, normalize: false })
        .then(() => {
          history.push(moveDetailsPath);
        })
        .catch(() => {
          setErrorMessage(`A server error occurred adding the shipment`);
        });
    } else if (isForServicesCounseling) {
      // routing and error handling handled in parent components
      submitHandler(updateMTOShipmentPayload);
    } else {
      submitHandler(updateMTOShipmentPayload)
        .then(() => {
          history.push(moveDetailsPath);
        })
        .catch(() => {
          setErrorMessage('A server error occurred editing the shipment details');
        });
    }
  };

  return (
    <Formik
      initialValues={initialValues}
      validateOnMount
      validateOnBlur
      validationSchema={schema}
      onSubmit={submitMTOShipment}
    >
      {({ values, isValid, isSubmitting, setValues, handleSubmit }) => {
        const { hasDeliveryAddress } = values;

        const handleUseCurrentResidenceChange = (e) => {
          const { checked } = e.target;
          if (checked) {
            // use current residence
            setValues({
              ...values,
              pickup: {
                ...values.pickup,
                address: currentResidence,
              },
            });
          } else {
            // Revert address
            setValues({
              ...values,
              pickup: {
                ...values.pickup,
                address: {
                  streetAddress1: '',
                  streetAddress2: '',
                  city: '',
                  state: '',
                  postalCode: '',
                },
              },
            });
          }
        };

        return (
          <>
            {isCancelModalVisible && (
              <SCRequestShipmentCancellationModal
                shipmentID={mtoShipment.id}
                onClose={setIsCancelModalVisible}
                onSubmit={handleDeleteShipment}
              />
            )}
            {errorMessage && (
              <Alert type="error" heading="An error occurred">
                {errorMessage}
              </Alert>
            )}
            {isTOO && mtoShipment.usesExternalVendor && (
              <Alert type="warning">
                The GHC prime contractor is not handling the shipment. Information will not be automatically shared with
                the movers handling it.
              </Alert>
            )}

            <div className={styles.ShipmentForm}>
              <div className={styles.headerWrapper}>
                <div>
                  <ShipmentTag shipmentType={shipmentType} shipmentNumber={shipmentNumber} />

                  <h1>{isCreatePage ? 'Add' : 'Edit'} shipment details</h1>
                </div>
                {!isCreatePage && (
                  <Button
                    type="button"
                    onClick={() => {
                      handleShowCancellationModal();
                    }}
                    unstyled
                  >
                    Delete shipment
                  </Button>
                )}
              </div>

              <SectionWrapper className={styles.weightAllowance}>
                <p>
                  <strong>Weight allowance: </strong>
                  {formatWeight(serviceMember.weightAllotment.totalWeightSelf)}
                </p>
              </SectionWrapper>

              <Form className={formStyles.form}>
                {isTOO && !isHHG && <ShipmentVendor />}

                {isNTSR && <ShipmentWeightInput userRole={userRole} />}

                {showPickupFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h2 className={styles.SectionHeaderExtraSpacing}>Pickup details</h2>
                    <Fieldset>
                      <DatePickerInput
                        name="pickup.requestedDate"
                        label="Requested pickup date"
                        id="requestedPickupDate"
                        validate={validateDate}
                      />
                    </Fieldset>

                    <AddressFields
                      name="pickup.address"
                      legend="Pickup location"
                      render={(fields) => (
                        <>
                          <Checkbox
                            data-testid="useCurrentResidence"
                            label="Use current address"
                            name="useCurrentResidence"
                            onChange={handleUseCurrentResidenceChange}
                            id="useCurrentResidenceCheckbox"
                          />
                          {fields}
                        </>
                      )}
                    />

                    <ContactInfoFields
                      name="pickup.agent"
                      legend={<div className={formStyles.legendContent}>Releasing agent {optionalLabel}</div>}
                      render={(fields) => {
                        return fields;
                      }}
                    />
                  </SectionWrapper>
                )}

                {isTOO && (isNTS || isNTSR) && (
                  <>
                    <StorageFacilityInfo userRole={userRole} />
                    <StorageFacilityAddress />
                  </>
                )}

                {isServiceCounselor && isNTSR && (
                  <>
                    <StorageFacilityInfo userRole={userRole} />
                    <StorageFacilityAddress />
                  </>
                )}

                {showDeliveryFields && (
                  <SectionWrapper className={formStyles.formSection}>
                    <h2 className={styles.SectionHeaderExtraSpacing}>Delivery details</h2>
                    <Fieldset>
                      <DatePickerInput
                        name="delivery.requestedDate"
                        label="Requested delivery date"
                        id="requestedDeliveryDate"
                        validate={validateDate}
                      />
                    </Fieldset>

                    {isNTSR ? (
                      <Fieldset legend="Delivery location">
                        <AddressFields
                          name="delivery.address"
                          render={(fields) => {
                            return fields;
                          }}
                        />
                        {displayDestinationType && (
                          <DropdownInput
                            label="Destination type"
                            name="destinationType"
                            options={shipmentDestinationAddressOptions}
                            id="destinationType"
                          />
                        )}
                      </Fieldset>
                    ) : (
                      <Fieldset legend="Delivery location">
                        <FormGroup>
                          <p>Does the customer know their delivery address yet?</p>
                          <div className={formStyles.radioGroup}>
                            <Field
                              as={Radio}
                              id="has-delivery-address"
                              label="Yes"
                              name="hasDeliveryAddress"
                              value="yes"
                              title="Yes, I know my delivery address"
                              checked={hasDeliveryAddress === 'yes'}
                            />
                            <Field
                              as={Radio}
                              id="no-delivery-address"
                              label="No"
                              name="hasDeliveryAddress"
                              value="no"
                              title="No, I do not know my delivery address"
                              checked={hasDeliveryAddress === 'no'}
                            />
                          </div>
                        </FormGroup>
                        {hasDeliveryAddress === 'yes' ? (
                          <>
                            <AddressFields
                              name="delivery.address"
                              render={(fields) => {
                                return fields;
                              }}
                            />
                            {displayDestinationType && (
                              <DropdownInput
                                label="Destination type"
                                name="destinationType"
                                options={shipmentDestinationAddressOptions}
                                id="destinationType"
                              />
                            )}
                          </>
                        ) : (
                          <p>
                            We can use the zip of their{' '}
                            {displayDestinationType ? 'HOR, HOS or PLEAD:' : 'new duty location:'}
                            <br />
                            <strong>
                              {newDutyStationAddress.city}, {newDutyStationAddress.state}{' '}
                              {newDutyStationAddress.postalCode}{' '}
                            </strong>
                          </p>
                        )}
                      </Fieldset>
                    )}

                    <ContactInfoFields
                      name="delivery.agent"
                      legend={<div className={formStyles.legendContent}>Receiving agent {optionalLabel}</div>}
                      render={(fields) => {
                        return fields;
                      }}
                    />
                  </SectionWrapper>
                )}

                <ShipmentFormRemarks
                  userRole={userRole}
                  customerRemarks={mtoShipment.customerRemarks}
                  counselorRemarks={mtoShipment.counselorRemarks}
                />

                {showAccountingCodes && (
                  <ShipmentAccountingCodes
                    TACs={TACs}
                    SACs={SACs}
                    onEditCodesClick={() => history.push(editOrdersPath)}
                    optional={isServiceCounselor}
                  />
                )}

                <div className={`${formStyles.formActions} ${styles.buttonGroup}`}>
                  <Button
                    data-testid="submitForm"
                    disabled={isSubmitting || !isValid}
                    type="submit"
                    onClick={handleSubmit}
                  >
                    Save
                  </Button>
                  <Button
                    type="button"
                    secondary
                    onClick={() => {
                      history.push(moveDetailsPath);
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </Form>
            </div>
          </>
        );
      }}
    </Formik>
  );
};

ShipmentForm.propTypes = {
  match: MatchShape,
  history: shape({
    push: func.isRequired,
  }),
  submitHandler: func.isRequired,
  isCreatePage: bool,
  isForServicesCounseling: bool,
  currentResidence: AddressShape.isRequired,
  newDutyStationAddress: SimpleAddressShape,
  selectedMoveType: string.isRequired,
  mtoShipment: HhgShipmentShape,
  moveTaskOrderID: string.isRequired,
  mtoShipments: arrayOf(MtoShipmentShape).isRequired,
  serviceMember: shape({
    weightAllotment: shape({
      totalWeightSelf: number,
    }),
  }).isRequired,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  userRole: oneOf(officeRoles).isRequired,
  displayDestinationType: bool,
};

ShipmentForm.defaultProps = {
  isCreatePage: false,
  isForServicesCounseling: false,
  match: { isExact: false, params: { moveCode: '', shipmentId: '' } },
  history: { push: () => {} },
  newDutyStationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
  mtoShipment: {
    id: '',
    customerRemarks: '',
    counselorRemarks: '',
    requestedPickupDate: '',
    requestedDeliveryDate: '',
    destinationAddress: {
      city: '',
      postalCode: '',
      state: '',
      streetAddress1: '',
    },
  },
  TACs: {},
  SACs: {},
  displayDestinationType: false,
};

export default ShipmentForm;
