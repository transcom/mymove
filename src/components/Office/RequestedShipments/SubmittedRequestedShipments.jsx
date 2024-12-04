import React, { useCallback, useState, useEffect } from 'react';
import { useFormik } from 'formik';
import * as PropTypes from 'prop-types';
import { Button, Checkbox, Fieldset } from '@trussworks/react-uswds';
import { generatePath, useParams, useNavigate } from 'react-router-dom';
import { debounce } from 'lodash';
import { connect } from 'react-redux';

import styles from './RequestedShipments.module.scss';

import { hasCounseling, hasMoveManagement } from 'utils/serviceItems';
import { isPPMOnly } from 'utils/shipments';
import ShipmentApprovalPreview from 'components/Office/ShipmentApprovalPreview/ShipmentApprovalPreview';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { tooRoutes } from 'constants/routes';
import { shipmentDestinationTypes } from 'constants/shipments';
import { permissionTypes } from 'constants/permissions';
import Restricted from 'components/Restricted/Restricted';
import { serviceItemCodes } from 'content/serviceItems';
import { shipmentTypeLabels } from 'content/shipments';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';
import { MoveTaskOrderShape, MTOServiceItemShape, OrdersInfoShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';
import { fieldValidationShape } from 'utils/displayFlags';
import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';
import { SHIPMENT_OPTIONS_URL, FEATURE_FLAG_KEYS } from 'shared/constants';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

// nts defaults show preferred pickup date and pickup address, flagged items when collapsed
// ntsr defaults shows preferred delivery date, storage facility address, delivery address, flagged items when collapsed
// Different things show when collapsed depending on if the shipment is an external vendor or not.
const showWhenCollapsedWithExternalVendor = {
  HHG_INTO_NTS_DOMESTIC: ['serviceOrderNumber', 'requestedDeliveryDate'],
  HHG_OUTOF_NTS_DOMESTIC: ['serviceOrderNumber', 'requestedPickupDate'],
};

const showWhenCollapsedWithGHCPrime = {
  HHG_INTO_NTS_DOMESTIC: ['tacType', 'requestedDeliveryDate'],
  HHG_OUTOF_NTS_DOMESTIC: ['ntsRecordedWeight', 'serviceOrderNumber', 'tacType', 'requestedPickupDate'],
};

const SubmittedRequestedShipments = ({
  mtoShipments,
  closeoutOffice,
  moveTaskOrder,
  allowancesInfo,
  ordersInfo,
  customerInfo,
  approveMTO,
  approveMTOShipment,
  handleAfterSuccess,
  missingRequiredOrdersInfo,
  errorIfMissing,
  displayDestinationType,
  mtoServiceItems,
  isMoveLocked,
  setFlashMessage,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [filteredShipments, setFilteredShipments] = useState([]);
  const [enableBoat, setEnableBoat] = useState(false);
  const [enableMobileHome, setEnableMobileHome] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setEnableBoat(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BOAT));
      setEnableMobileHome(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.MOBILE_HOME));
    };
    fetchData();
  }, []);

  const filterPrimeShipments = mtoShipments.filter((shipment) => !shipment.usesExternalVendor);

  const filterShipments = (formikShipmentIds) => {
    return mtoShipments.filter(({ id }) => formikShipmentIds.includes(id));
  };

  const ordersLOA = {
    tac: ordersInfo.tacMDC,
    sac: ordersInfo.sacSDN,
    ntsTac: ordersInfo.NTStac,
    ntsSac: ordersInfo.NTSsac,
  };

  const { moveCode } = useParams();
  const navigate = useNavigate();
  const handleButtonDropdownChange = (e) => {
    const selectedOption = e.target.value;

    const addShipmentPath = `${generatePath(tooRoutes.SHIPMENT_ADD_PATH, {
      moveCode,
      shipmentType: selectedOption,
    })}`;

    navigate(addShipmentPath);
  };

  const shipmentDisplayInfo = (shipment, dutyLocationPostal) => {
    const destType = displayDestinationType ? shipmentDestinationTypes[shipment.destinationType] : null;

    return {
      ...shipment,
      heading: shipmentTypeLabels[shipment.shipmentType],
      isDiversion: shipment.diversion,
      shipmentStatus: shipment.status,
      destinationAddress: shipment.destinationAddress || dutyLocationPostal,
      destinationType: destType,
      displayDestinationType,
      closeoutOffice,
    };
  };

  const allowedShipmentOptions = () => {
    return (
      <>
        <option data-testid="hhgOption" value={SHIPMENT_OPTIONS_URL.HHG}>
          HHG
        </option>
        <option value={SHIPMENT_OPTIONS_URL.PPM}>PPM</option>
        <option value={SHIPMENT_OPTIONS_URL.NTS}>NTS</option>
        <option value={SHIPMENT_OPTIONS_URL.NTSrelease}>NTS-release</option>
        {enableBoat && <option value={SHIPMENT_OPTIONS_URL.BOAT}>Boat</option>}
        {enableMobileHome && <option value={SHIPMENT_OPTIONS_URL.MOBILE_HOME}>Mobile Home</option>}
      </>
    );
  };

  const formik = useFormik({
    initialValues: {
      shipmentManagementFee: true,
      counselingFee: false,
      shipments: [],
    },
    onSubmit: (values, { setSubmitting }) => {
      const mtoApprovalServiceItemCodes = {
        serviceCodeMS: values.shipmentManagementFee && !moveTaskOrder.availableToPrimeAt,
        serviceCodeCS: values.counselingFee,
      };

      approveMTO(
        {
          moveTaskOrderID: moveTaskOrder.id,
          ifMatchETag: moveTaskOrder.eTag,
          mtoApprovalServiceItemCodes,
          normalize: false,
        },
        {
          onSuccess: async () => {
            try {
              await Promise.all(
                filteredShipments.map((shipment) => {
                  let operationPath = 'shipment.approveShipment';

                  if (shipment.approvedDate && moveTaskOrder.availableToPrimeAt) {
                    operationPath = 'shipment.approveShipmentDiversion';
                  }
                  return approveMTOShipment(
                    {
                      shipmentID: shipment.id,
                      operationPath,
                      ifMatchETag: shipment.eTag,
                      normalize: false,
                    },
                    {
                      onError: () => {
                        // TODO: Decide if we want to display an error notice, log error event, or retry
                        setSubmitting(false);
                        setFlashMessage(null);
                      },
                    },
                  );
                }),
              );
              setFlashMessage('TASK_ORDER_CREATE_SUCCESS', 'success', 'Task order created successfully.');
              handleAfterSuccess('../mto', { showMTOpostedMessage: true });
            } catch {
              setSubmitting(false);
            }
          },
          onError: () => {
            // TODO: Decide if we want to display an error notice, log error event, or retry
            setSubmitting(false);
          },
        },
      );
      //
    },
  });

  const handleReviewClick = () => {
    setFilteredShipments(filterShipments(formik.values.shipments));
    setIsModalVisible(true);
  };

  // if showing service items on a move with Prime shipments, enable button when shipment and service item are selected and there is no missing required Orders information
  // if not showing service items on a move with Prime shipments, enable button if a shipment is selected and there is no missing required Orders information
  const primeShipmentsForApproval = moveTaskOrder.availableToPrimeAt
    ? formik.values.shipments.length > 0 && !missingRequiredOrdersInfo
    : formik.values.shipments.length > 0 &&
      (formik.values.counselingFee || formik.values.shipmentManagementFee) &&
      !missingRequiredOrdersInfo;

  // on a move with only External Vendor shipments enable button if a service item is selected
  const externalVendorShipmentsOnly = formik.values.counselingFee || formik.values.shipmentManagementFee;

  // Check that there are Prime-handled shipments before determining if the button should be enabled
  const isButtonEnabled = filterPrimeShipments.length > 0 ? primeShipmentsForApproval : externalVendorShipmentsOnly;

  const dutyLocationPostal = { postalCode: ordersInfo.newDutyLocation?.address?.postalCode };

  // Hide counseling line item if prime counseling is already in the service items or if service counseling has been applied
  const hideCounselingCheckbox = hasCounseling(mtoServiceItems) || moveTaskOrder?.serviceCounselingCompletedAt;

  // Hide move management line item if it is already in the service items or for PPM only moves
  const hideMoveManagementCheckbox = hasMoveManagement(mtoServiceItems) || isPPMOnly(mtoShipments);

  // Disable move management checkbox
  const moveManagementDisabled = true;

  // Check the move management box
  const moveManagementChecked = true;

  // If we are hiding both counseling and move management then hide the entire service item form
  const hideAddServiceItemsForm = hideCounselingCheckbox && hideMoveManagementCheckbox;

  // Adding an inline function will break the debounce fix and allow multiple submits
  // RA Validator Status: RA Accepted
  // eslint-disable-next-line react-hooks/exhaustive-deps
  const debouncedSubmit = useCallback(debounce(formik.handleSubmit, 5000, { leading: true }), []);

  return (
    <div className={styles.RequestedShipments} data-testid="requested-shipments">
      <div
        id="approvalConfirmationModal"
        data-testid="approvalConfirmationModal"
        style={{ display: isModalVisible ? 'block' : 'none' }}
      >
        <ShipmentApprovalPreview
          mtoShipments={filteredShipments}
          ordersInfo={ordersInfo}
          allowancesInfo={allowancesInfo}
          customerInfo={customerInfo}
          setIsModalVisible={setIsModalVisible}
          onSubmit={debouncedSubmit}
          counselingFee={formik.values.counselingFee}
          shipmentManagementFee={formik.values.shipmentManagementFee}
        />
      </div>

      <form onSubmit={formik.handleSubmit}>
        <div className={styles.sectionHeader}>
          <h2>Requested shipments</h2>
          <div className={styles.buttonDropdown}>
            {!isMoveLocked && (
              <Restricted to={permissionTypes.createTxoShipment}>
                <ButtonDropdown
                  ariaLabel="Add a new shipment"
                  data-testid="addShipmentButton"
                  onChange={handleButtonDropdownChange}
                >
                  <option value="" label="Add a new shipment">
                    Add a new shipment
                  </option>
                  {allowedShipmentOptions()}
                </ButtonDropdown>
              </Restricted>
            )}
          </div>
        </div>
        <div className={shipmentCardsStyles.shipmentCards}>
          {mtoShipments &&
            mtoShipments.map((shipment) => {
              const editUrl = `../${generatePath(tooRoutes.SHIPMENT_EDIT_PATH, {
                shipmentId: shipment.id,
              })}`;

              return (
                <ShipmentDisplay
                  key={shipment.id}
                  isSubmitted
                  shipmentId={shipment.id}
                  shipmentType={shipment.shipmentType}
                  displayInfo={shipmentDisplayInfo(shipment, dutyLocationPostal)}
                  ordersLOA={ordersLOA}
                  errorIfMissing={errorIfMissing[shipment.shipmentType]}
                  showWhenCollapsed={
                    shipment.usesExternalVendor
                      ? showWhenCollapsedWithExternalVendor[shipment.shipmentType]
                      : showWhenCollapsedWithGHCPrime[shipment.shipmentType]
                  }
                  editURL={editUrl}
                  /* eslint-disable-next-line react/jsx-props-no-spreading */
                  {...formik.getFieldProps(`shipments`)}
                  isMoveLocked={isMoveLocked}
                />
              );
            })}
        </div>

        <Restricted to={permissionTypes.updateShipment}>
          <div className={styles.serviceItems}>
            {!hideAddServiceItemsForm && (
              <>
                <h2>Add service items to this move</h2>
                <Fieldset legend="MTO service items" legendsronly="true" id="input-type-fieldset">
                  {!hideMoveManagementCheckbox && (
                    <Checkbox
                      id="shipmentManagementFee"
                      label={serviceItemCodes.MS}
                      name="shipmentManagementFee"
                      onChange={formik.handleChange}
                      checked={moveManagementChecked}
                      disabled={moveManagementDisabled}
                      data-testid="shipmentManagementFee"
                    />
                  )}
                  {hideCounselingCheckbox ? (
                    <p className={styles.serviceCounselingCompleted} data-testid="services-counseling-completed-text">
                      The customer has received counseling for this move.
                    </p>
                  ) : (
                    <Checkbox
                      id="counselingFee"
                      label={serviceItemCodes.CS}
                      name="counselingFee"
                      onChange={formik.handleChange}
                      data-testid="counselingFee"
                      disabled={isMoveLocked}
                    />
                  )}
                </Fieldset>
              </>
            )}
            <Button
              data-testid="shipmentApproveButton"
              className={styles.approveButton}
              onClick={handleReviewClick}
              type="button"
              disabled={!isButtonEnabled || isMoveLocked}
            >
              <span>Approve selected</span>
            </Button>
          </div>
        </Restricted>
      </form>
    </div>
  );
};

SubmittedRequestedShipments.propTypes = {
  mtoShipments: PropTypes.arrayOf(ShipmentShape).isRequired,
  ordersInfo: OrdersInfoShape.isRequired,
  allowancesInfo: PropTypes.shape({
    branch: PropTypes.string,
    grade: PropTypes.string,
    totalWeight: PropTypes.number,
    progear: PropTypes.number,
    spouseProgear: PropTypes.number,
    storageInTransit: PropTypes.number,
    dependents: PropTypes.bool,
  }).isRequired,
  customerInfo: PropTypes.shape({
    name: PropTypes.string,
    dodId: PropTypes.string,
    phone: PropTypes.string,
    email: PropTypes.string,
    currentAddress: PropTypes.shape({
      streetAddress1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postalCode: PropTypes.string,
    }),
    backupContactName: PropTypes.string,
    backupContactPhone: PropTypes.string,
    backupContactEmail: PropTypes.string,
  }).isRequired,
  approveMTO: PropTypes.func,
  approveMTOShipment: PropTypes.func,
  moveTaskOrder: MoveTaskOrderShape,
  missingRequiredOrdersInfo: PropTypes.bool,
  handleAfterSuccess: PropTypes.func,
  errorIfMissing: PropTypes.objectOf(PropTypes.arrayOf(fieldValidationShape)),
  displayDestinationType: PropTypes.bool,
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
};

SubmittedRequestedShipments.defaultProps = {
  moveTaskOrder: {},
  approveMTO: () => Promise.resolve(),
  approveMTOShipment: () => Promise.resolve(),
  missingRequiredOrdersInfo: false,
  handleAfterSuccess: () => {},
  errorIfMissing: {},
  displayDestinationType: false,
  mtoServiceItems: [],
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
};
export default connect(() => ({}), mapDispatchToProps)(SubmittedRequestedShipments);
