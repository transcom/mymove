import React, { useState } from 'react';
import { arrayOf, bool, func, number, shape, string, oneOf } from 'prop-types';
import { Field, Formik } from 'formik';
import { generatePath } from 'react-router';
import { useQueryClient, useMutation } from '@tanstack/react-query';
import { Alert, Button, Checkbox, Fieldset, FormGroup, Radio } from '@trussworks/react-uswds';

import getShipmentOptions from '../../Customer/MtoShipmentForm/getShipmentOptions';
import { CloseoutOfficeInput } from '../../form/fields/CloseoutOfficeInput';

import styles from './ShipmentForm.module.scss';
import ppmShipmentSchema from './ppmShipmentSchema';

import SITCostDetails from 'components/Office/SITCostDetails/SITCostDetails';
import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { ContactInfoFields } from 'components/form/ContactInfoFields/ContactInfoFields';
import { DatePickerInput, DropdownInput } from 'components/form/fields';
import { Form } from 'components/form/Form';
import DestinationZIPInfo from 'components/Office/DestinationZIPInfo/DestinationZIPInfo';
import OriginZIPInfo from 'components/Office/OriginZIPInfo/OriginZIPInfo';
import ShipmentAccountingCodes from 'components/Office/ShipmentAccountingCodes/ShipmentAccountingCodes';
import ShipmentCustomerSIT from 'components/Office/ShipmentCustomerSIT/ShipmentCustomerSIT';
import ShipmentFormRemarks from 'components/Office/ShipmentFormRemarks/ShipmentFormRemarks';
import ShipmentIncentiveAdvance from 'components/Office/ShipmentIncentiveAdvance/ShipmentIncentiveAdvance';
import ShipmentVendor from 'components/Office/ShipmentVendor/ShipmentVendor';
import ShipmentWeight from 'components/Office/ShipmentWeight/ShipmentWeight';
import ShipmentWeightInput from 'components/Office/ShipmentWeightInput/ShipmentWeightInput';
import StorageFacilityAddress from 'components/Office/StorageFacilityAddress/StorageFacilityAddress';
import StorageFacilityInfo from 'components/Office/StorageFacilityInfo/StorageFacilityInfo';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { MOVES, MTO_SHIPMENTS } from 'constants/queryKeys';
import { servicesCounselingRoutes, tooRoutes } from 'constants/routes';
import { shipmentDestinationTypes } from 'constants/shipments';
import { officeRoles, roleTypes } from 'constants/userRoles';
import { deleteShipment, updateMoveCloseoutOffice } from 'services/ghcApi';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import formStyles from 'styles/form.module.scss';
import { AccountingCodesShape } from 'types/accountingCodes';
import { AddressShape, SimpleAddressShape } from 'types/address';
import { MatchShape } from 'types/officeShapes';
import { ShipmentShape } from 'types/shipment';
import { TransportationOfficeShape } from 'types/transportationOffice';
import {
  formatMtoShipmentForAPI,
  formatMtoShipmentForDisplay,
  formatPpmShipmentForAPI,
  formatPpmShipmentForDisplay,
} from 'utils/formatMtoShipment';
import { formatWeight, dropdownInputOptions } from 'utils/formatters';
import { validateDate, validatePostalCode } from 'utils/validation';

const ShipmentForm = (props) => {
  const {
    match,
    history,
    originDutyLocationAddress,
    newDutyLocationAddress,
    shipmentType,
    isCreatePage,
    isForServicesCounseling,
    mtoShipment,
    submitHandler,
    onUpdate,
    mtoShipments,
    serviceMember,
    currentResidence,
    moveTaskOrderID,
    TACs,
    SACs,
    userRole,
    displayDestinationType,
    isAdvancePage,
    move,
  } = props;

  const { moveCode } = match.params;
  const [errorMessage, setErrorMessage] = useState(null);
  const [isCancelModalVisible, setIsCancelModalVisible] = useState(false);

  const shipments = mtoShipments;

  const queryClient = useQueryClient();
  const { mutate: mutateMTOShipmentStatus } = useMutation(deleteShipment, {
    onSuccess: (_, variables) => {
      const updatedMTOShipment = mtoShipment;
      // Update mtoShipments with our updated status and set query data to match
      shipments[mtoShipments.findIndex((shipment) => shipment.id === updatedMTOShipment.id)] = updatedMTOShipment;
      queryClient.setQueryData([MTO_SHIPMENTS, updatedMTOShipment.moveTaskOrderID, false], mtoShipments);
      // InvalidateQuery tells other components using this data that they need to re-fetch
      // This allows the requestCancellation button to update immediately
      queryClient.invalidateQueries([MTO_SHIPMENTS, variables.moveTaskOrderID]);

      history.goBack();
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      setErrorMessage(errorMsg);
    },
  });

  const { mutate: mutateMoveCloseoutOffice } = useMutation(updateMoveCloseoutOffice, {
    onSuccess: () => {
      queryClient.invalidateQueries([MOVES, moveCode]);
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

  const isHHG = shipmentType === SHIPMENT_OPTIONS.HHG;
  const isNTS = shipmentType === SHIPMENT_OPTIONS.NTS;
  const isNTSR = shipmentType === SHIPMENT_OPTIONS.NTSR;
  const isPPM = shipmentType === SHIPMENT_OPTIONS.PPM;

  const showAccountingCodes = isNTS || isNTSR;

  const isTOO = userRole === roleTypes.TOO;
  const isServiceCounselor = userRole === roleTypes.SERVICES_COUNSELOR;
  const showCloseoutOffice =
    isServiceCounselor && isPPM && (serviceMember.agency === 'ARMY' || serviceMember.agency === 'AIR_FORCE');

  const shipmentDestinationAddressOptions = dropdownInputOptions(shipmentDestinationTypes);

  const shipmentNumber = isHHG ? getShipmentNumber() : null;
  const initialValues = isPPM
    ? formatPpmShipmentForDisplay(
        isCreatePage
          ? { closeoutOffice: move.closeoutOffice }
          : {
              counselorRemarks: mtoShipment.counselorRemarks,
              ppmShipment: mtoShipment.ppmShipment,
              closeoutOffice: move.closeoutOffice,
            },
      )
    : formatMtoShipmentForDisplay(
        isCreatePage
          ? { userRole, shipmentType }
          : { userRole, shipmentType, agents: mtoShipment.mtoAgents, ...mtoShipment },
      );

  let showDeliveryFields;
  let showPickupFields;
  let schema;

  if (isPPM) {
    schema = ppmShipmentSchema({
      estimatedIncentive: initialValues.estimatedIncentive || 0,
      weightAllotment: serviceMember.weightAllotment,
      advanceAmountRequested: mtoShipment.ppmShipment?.advanceAmountRequested,
      hasRequestedAdvance: mtoShipment.ppmShipment?.hasRequestedAdvance,
      isAdvancePage,
      showCloseoutOffice,
    });
  } else {
    const shipmentOptions = getShipmentOptions(shipmentType, userRole);

    showDeliveryFields = shipmentOptions.showDeliveryFields;
    showPickupFields = shipmentOptions.showPickupFields;
    schema = shipmentOptions.schema;
  }

  const optionalLabel = <span className={formStyles.optional}>Optional</span>;

  const moveDetailsRoute = isTOO ? tooRoutes.MOVE_VIEW_PATH : servicesCounselingRoutes.MOVE_VIEW_PATH;
  const moveDetailsPath = generatePath(moveDetailsRoute, { moveCode });

  const editOrdersRoute = isTOO ? tooRoutes.ORDERS_EDIT_PATH : servicesCounselingRoutes.ORDERS_EDIT_PATH;
  const editOrdersPath = generatePath(editOrdersRoute, { moveCode });

  const submitMTOShipment = (formValues, actions) => {
    //* PPM Shipment *//
    if (isPPM) {
      const ppmShipmentBody = formatPpmShipmentForAPI(formValues);
      // Add a PPM shipment
      if (isCreatePage) {
        const body = { ...ppmShipmentBody, moveTaskOrderID };
        submitHandler(
          { body, normalize: false },
          {
            onSuccess: (newMTOShipment) => {
              const currentPath = generatePath(servicesCounselingRoutes.SHIPMENT_EDIT_PATH, {
                moveCode,
                shipmentId: newMTOShipment.id,
              });
              const advancePath = generatePath(servicesCounselingRoutes.SHIPMENT_ADVANCE_PATH, {
                moveCode,
                shipmentId: newMTOShipment.id,
              });
              if (formValues.closeoutOffice.id) {
                mutateMoveCloseoutOffice({
                  locator: moveCode,
                  ifMatchETag: move.eTag,
                  body: { closeoutOfficeId: formValues.closeoutOffice.id },
                });
              }
              history.replace(currentPath);
              history.push(advancePath);
            },
            onError: () => {
              actions.setSubmitting(false);
              setErrorMessage(`A server error occurred adding the shipment`);
            },
          },
        );
        return;
      }
      // Edit a PPM Shipment
      const updatePPMPayload = {
        moveTaskOrderID,
        shipmentID: mtoShipment.id,
        ifMatchETag: mtoShipment.eTag,
        normalize: false,
        body: ppmShipmentBody,
        locator: move.locator,
        moveETag: move.eTag,
      };

      submitHandler(updatePPMPayload, {
        onSuccess: () => {
          if (!isAdvancePage && formValues.closeoutOffice.id) {
            // if we have a closeout office, we must be on the first page of creating a PPM shipment,
            // as a SC so we should update the closeout office and redirect to the advance page
            mutateMoveCloseoutOffice(
              {
                locator: moveCode,
                ifMatchETag: move.eTag,
                body: { closeoutOfficeId: formValues.closeoutOffice.id },
              },
              {
                onSuccess: () => {
                  const advancePath = generatePath(servicesCounselingRoutes.SHIPMENT_ADVANCE_PATH, {
                    moveCode,
                    shipmentId: mtoShipment.id,
                  });
                  actions.setSubmitting(false);
                  history.push(advancePath);
                  onUpdate('success');
                },
                onError: () => {
                  history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
                  onUpdate('error');
                },
              },
            );
          } else {
            // if we don't have a closeout office, we're either on the advance page for a PPM as a SC or the first page of a PPM as a TOO.
            // In any case, we're done now and can head back to the move viewÎ
            history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
            onUpdate('success');
          }
        },
      });
      return;
    }

    //* MTO Shipments *//

    const {
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
    } = formValues;

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
      shipmentType,
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

    // Add a MTO Shipment (only a Service Counselor can add a shipment)
    if (isCreatePage) {
      const body = { ...pendingMtoShipment, moveTaskOrderID };
      submitHandler(
        { body, normalize: false },
        {
          onSuccess: () => {
            history.push(moveDetailsPath);
          },
          onError: () => {
            setErrorMessage(`A server error occurred adding the shipment`);
          },
        },
      );
    }
    // Edit MTO as Service Counselor
    else if (isForServicesCounseling) {
      // error handling handled in parent components
      submitHandler(updateMTOShipmentPayload, {
        onSuccess: () => {
          history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
          onUpdate('success');
        },
        onError: () => {
          history.push(generatePath(servicesCounselingRoutes.MOVE_VIEW_PATH, { moveCode }));
          onUpdate('error');
        },
      });
    }
    // Edit a MTO Shipment as TOO
    else {
      submitHandler(updateMTOShipmentPayload, {
        onSuccess: () => {
          history.push(moveDetailsPath);
        },
        onError: () => {
          setErrorMessage('A server error occurred editing the shipment details');
        },
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
      {({ values, isValid, isSubmitting, setValues, handleSubmit, errors }) => {
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
            <ConnectedDestructiveShipmentConfirmationModal
              isOpen={isCancelModalVisible}
              shipmentID={mtoShipment.id}
              onClose={setIsCancelModalVisible}
              onSubmit={handleDeleteShipment}
            />
            {errorMessage && (
              <Alert type="error" headingLevel="h4" heading="An error occurred">
                {errorMessage}
              </Alert>
            )}
            {isTOO && mtoShipment.usesExternalVendor && (
              <Alert headingLevel="h4" type="warning">
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
                              {newDutyLocationAddress.city}, {newDutyLocationAddress.state}{' '}
                              {newDutyLocationAddress.postalCode}{' '}
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

                {isPPM && !isAdvancePage && (
                  <>
                    <OriginZIPInfo
                      postalCodeValidator={validatePostalCode}
                      currentZip={originDutyLocationAddress.postalCode}
                    />
                    <DestinationZIPInfo
                      postalCodeValidator={validatePostalCode}
                      dutyZip={newDutyLocationAddress.postalCode}
                    />
                    {showCloseoutOffice && (
                      <SectionWrapper>
                        <h2>Closeout office</h2>
                        <CloseoutOfficeInput
                          hint="If there is more than one PPM for this move, the closeout office will be the same for all your PPMs."
                          name="closeoutOffice"
                          placeholder="Start typing a closeout location..."
                          label="Closeout location"
                          displayAddress
                        />
                      </SectionWrapper>
                    )}
                    <ShipmentCustomerSIT />
                    <ShipmentWeight authorizedWeight={serviceMember.weightAllotment.totalWeightSelf.toString()} />
                  </>
                )}

                {isPPM && isAdvancePage && isServiceCounselor && mtoShipment.ppmShipment?.sitExpected && (
                  <SITCostDetails
                    cost={mtoShipment.ppmShipment?.sitEstimatedCost}
                    weight={mtoShipment.ppmShipment?.sitEstimatedWeight}
                    sitLocation={mtoShipment.ppmShipment?.sitLocation}
                    originZip={mtoShipment.ppmShipment?.pickupPostalCode}
                    destinationZip={mtoShipment.ppmShipment?.destinationPostalCode}
                    departureDate={mtoShipment.ppmShipment?.sitEstimatedDepartureDate}
                    entryDate={mtoShipment.ppmShipment?.sitEstimatedEntryDate}
                  />
                )}

                {isPPM && isAdvancePage && (
                  <ShipmentIncentiveAdvance
                    values={values}
                    estimatedIncentive={mtoShipment.ppmShipment?.estimatedIncentive}
                  />
                )}

                {(!isPPM || (isPPM && isAdvancePage)) && (
                  <ShipmentFormRemarks
                    userRole={userRole}
                    shipmentType={shipmentType}
                    customerRemarks={mtoShipment.customerRemarks}
                    counselorRemarks={mtoShipment.counselorRemarks}
                    showHint={false}
                    error={
                      errors.counselorRemarks &&
                      (values.advanceRequested !== mtoShipment.ppmShipment?.hasRequestedAdvance ||
                        values.advance !== mtoShipment.ppmShipment?.advanceAmountRequested)
                    }
                  />
                )}

                {showAccountingCodes && (
                  <ShipmentAccountingCodes
                    TACs={TACs}
                    SACs={SACs}
                    onEditCodesClick={() => history.push(editOrdersPath)}
                    optional={isServiceCounselor}
                  />
                )}

                <div className={`${formStyles.formActions} ${styles.buttonGroup}`}>
                  {!isPPM && (
                    <Button
                      data-testid="submitForm"
                      disabled={isSubmitting || !isValid}
                      type="submit"
                      onClick={handleSubmit}
                    >
                      Save
                    </Button>
                  )}
                  <Button
                    type="button"
                    secondary
                    onClick={() => {
                      history.push(moveDetailsPath);
                    }}
                  >
                    Cancel
                  </Button>
                  {isPPM && (
                    <Button
                      data-testid="submitForm"
                      disabled={isSubmitting || !isValid}
                      type="submit"
                      onClick={handleSubmit}
                    >
                      Save and Continue
                    </Button>
                  )}
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
  onUpdate: func,
  isCreatePage: bool,
  isForServicesCounseling: bool,
  currentResidence: AddressShape.isRequired,
  originDutyLocationAddress: SimpleAddressShape,
  newDutyLocationAddress: SimpleAddressShape,
  shipmentType: string.isRequired,
  mtoShipment: ShipmentShape,
  moveTaskOrderID: string.isRequired,
  mtoShipments: arrayOf(ShipmentShape).isRequired,
  serviceMember: shape({
    weightAllotment: shape({
      totalWeightSelf: number,
    }),
    agency: string.isRequired,
  }).isRequired,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  userRole: oneOf(officeRoles).isRequired,
  displayDestinationType: bool,
  isAdvancePage: bool,
  move: shape({
    eTag: string,
    id: string,
    closeoutOffice: TransportationOfficeShape,
  }),
};

ShipmentForm.defaultProps = {
  isCreatePage: false,
  isForServicesCounseling: false,
  match: { isExact: false, params: { moveCode: '', shipmentId: '' } },
  history: { push: () => {} },
  onUpdate: () => {},
  originDutyLocationAddress: {
    city: '',
    state: '',
    postalCode: '',
  },
  newDutyLocationAddress: {
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
  isAdvancePage: false,
  move: {},
};

export default ShipmentForm;
