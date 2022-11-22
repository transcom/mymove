import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as PropTypes from 'prop-types';
import { Button, Checkbox, Fieldset } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import styles from './RequestedShipments.module.scss';

import ShipmentApprovalPreview from 'components/Office/ShipmentApprovalPreview/ShipmentApprovalPreview';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { tooRoutes } from 'constants/routes';
import { shipmentDestinationTypes } from 'constants/shipments';
import { permissionTypes } from 'constants/permissions';
import Restricted from 'components/Restricted/Restricted';
import { serviceItemCodes } from 'content/serviceItems';
import { shipmentTypeLabels } from 'content/shipments';
import shipmentCardsStyles from 'styles/shipmentCards.module.scss';
import { MoveTaskOrderShape, OrdersInfoShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';

// nts defaults show preferred pickup date and pickup address, flagged items when collapsed
// ntsr defaults shows preferred delivery date, storage facility address, destination address, flagged items when collapsed
// Different things show when collapsed depending on if the shipment is an external vendor or not.
const showWhenCollapsedWithExternalVendor = {
  HHG_INTO_NTS_DOMESTIC: ['serviceOrderNumber'],
  HHG_OUTOF_NTS_DOMESTIC: ['serviceOrderNumber'],
};

const showWhenCollapsedWithGHCPrime = {
  HHG_INTO_NTS_DOMESTIC: ['tacType'],
  HHG_OUTOF_NTS_DOMESTIC: ['ntsRecordedWeight', 'serviceOrderNumber', 'tacType'],
};

const SubmittedRequestedShipments = ({
  mtoShipments,
  moveTaskOrder,
  allowancesInfo,
  ordersInfo,
  customerInfo,
  approveMTO,
  approveMTOShipment,
  handleAfterSuccess,
  missingRequiredOrdersInfo,
  moveCode,
  errorIfMissing,
  displayDestinationType,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [filteredShipments, setFilteredShipments] = useState([]);

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
    };
  };

  const formik = useFormik({
    initialValues: {
      shipmentManagementFee: false,
      counselingFee: false,
      shipments: [],
    },
    onSubmit: (values, { setSubmitting }) => {
      const mtoApprovalServiceItemCodes = {
        serviceCodeMS: values.shipmentManagementFee,
        serviceCodeCS: values.counselingFee,
      };

      // The MTO has not yet been approved so resolve before updating the shipment statuses and creating accessorial service items
      if (!moveTaskOrder.availableToPrimeAt) {
        approveMTO({
          moveTaskOrderID: moveTaskOrder.id,
          ifMatchETag: moveTaskOrder.eTag,
          mtoApprovalServiceItemCodes,
          normalize: false,
        })
          .then(() => {
            Promise.all(
              filteredShipments.map((shipment) =>
                approveMTOShipment({
                  shipmentID: shipment.id,
                  operationPath: 'shipment.approveShipment',
                  ifMatchETag: shipment.eTag,
                  normalize: false,
                }),
              ),
            )
              .then(() => {
                handleAfterSuccess('mto', { showMTOpostedMessage: true });
              })
              .catch(() => {
                // TODO: Decide if we want to display an error notice, log error event, or retry
                setSubmitting(false);
              });
          })
          .catch(() => {
            // TODO: Decide if we want to display an error notice, log error event, or retry
            setSubmitting(false);
          });
      } else {
        // The MTO was previously approved along with at least one shipment, only update the new shipment statuses
        Promise.all(
          filteredShipments.map((shipment) => {
            let operationPath = 'shipment.approveShipment';

            if (shipment.approvedDate) {
              operationPath = 'shipment.approveShipmentDiversion';
            }

            return approveMTOShipment({
              shipmentID: shipment.id,
              operationPath,
              ifMatchETag: shipment.eTag,
              normalize: false,
            });
          }),
        )
          .then(() => {
            handleAfterSuccess('mto');
          })
          .catch(() => {
            // TODO: Decide if we want to display an error notice, log error event, or retry
            setSubmitting(false);
          });
      }
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

  // on a move with only External Vendor shipments enable button if a a service item is selected
  const externalVendorShipmentsOnly = formik.values.counselingFee || formik.values.shipmentManagementFee;

  // Check that there are Prime-handled shipments before determining if the button should be enabled
  const isButtonEnabled = filterPrimeShipments.length > 0 ? primeShipmentsForApproval : externalVendorShipmentsOnly;

  const dutyLocationPostal = { postalCode: ordersInfo.newDutyLocation?.address?.postalCode };

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
          onSubmit={formik.handleSubmit}
          counselingFee={formik.values.counselingFee}
          shipmentManagementFee={formik.values.shipmentManagementFee}
        />
      </div>

      <form onSubmit={formik.handleSubmit}>
        <h2>Requested shipments</h2>
        <div className={shipmentCardsStyles.shipmentCards}>
          {mtoShipments &&
            mtoShipments.map((shipment) => {
              const editUrl = generatePath(tooRoutes.SHIPMENT_EDIT_PATH, {
                moveCode,
                shipmentId: shipment.id,
              });

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
                />
              );
            })}
        </div>

        <Restricted to={permissionTypes.updateShipment}>
          <div className={styles.serviceItems}>
            {!moveTaskOrder.availableToPrimeAt && (
              <>
                <h2>Add service items to this move</h2>
                <Fieldset legend="MTO service items" legendsronly="true" id="input-type-fieldset">
                  <Checkbox
                    id="shipmentManagementFee"
                    label={serviceItemCodes.MS}
                    name="shipmentManagementFee"
                    onChange={formik.handleChange}
                  />
                  {moveTaskOrder.serviceCounselingCompletedAt ? (
                    <p className={styles.serviceCounselingCompleted} data-testid="services-counseling-completed-text">
                      The customer has received counseling for this move.
                    </p>
                  ) : (
                    <Checkbox
                      id="counselingFee"
                      label={serviceItemCodes.CS}
                      name="counselingFee"
                      onChange={formik.handleChange}
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
              disabled={!isButtonEnabled}
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
    rank: PropTypes.string,
    weightAllowance: PropTypes.number,
    authorizedWeight: PropTypes.number,
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
  moveCode: PropTypes.string.isRequired,
  handleAfterSuccess: PropTypes.func,
  errorIfMissing: PropTypes.shape({}),
  displayDestinationType: PropTypes.bool,
};

SubmittedRequestedShipments.defaultProps = {
  moveTaskOrder: {},
  approveMTO: () => Promise.resolve(),
  approveMTOShipment: () => Promise.resolve(),
  missingRequiredOrdersInfo: false,
  handleAfterSuccess: () => {},
  errorIfMissing: {},
  displayDestinationType: false,
};

export default SubmittedRequestedShipments;
