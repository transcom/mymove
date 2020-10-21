import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as PropTypes from 'prop-types';
import { Button, Checkbox, Fieldset } from '@trussworks/react-uswds';

import ShipmentApprovalPreview from '../ShipmentApprovalPreview';

import styles from './RequestedShipments.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import {
  MTOAgentShape,
  MTOShipmentShape,
  MoveTaskOrderShape,
  MTOServiceItemShape,
  OrdersInfoShape,
} from 'types/moveOrder';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import { ReactComponent as FormCheckmarkIcon } from 'shared/icon/form-checkmark.svg';
import { formatDateFromIso } from 'shared/formatters';

const RequestedShipments = ({
  mtoShipments,
  ordersInfo,
  allowancesInfo,
  customerInfo,
  mtoAgents,
  shipmentsStatus,
  mtoServiceItems,
  moveTaskOrder,
  approveMTO,
  approveMTOShipment,
}) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [filteredShipments, setFilteredShipments] = useState([]);

  const filterShipments = (formikShipmentIds) => {
    return mtoShipments.filter(({ id }) => formikShipmentIds.includes(id));
  };

  const formik = useFormik({
    initialValues: {
      shipmentManagementFee: false,
      counselingFee: false,
      shipments: [],
    },
    onSubmit: (values, { setSubmitting }) => {
      const requests = [
        Promise.all(
          filteredShipments.map((shipment) =>
            approveMTOShipment(moveTaskOrder.id, shipment.id, 'APPROVED', shipment.eTag),
          ),
        ),
      ];
      const mtoApprovalServiceItemCodes = {
        serviceCodeMS: values.shipmentManagementFee,
        serviceCodeCS: values.counselingFee,
      };

      // if mto is not yet approved, add request to approve it
      if (!moveTaskOrder.availableToPrimeAt) {
        requests.push(approveMTO(moveTaskOrder.id, moveTaskOrder.eTag, mtoApprovalServiceItemCodes));
      }

      Promise.all(requests)
        .then((results) => {
          if (results[0].every((shipmentResult) => shipmentResult.response.status === 200)) {
            if (results[1]?.response?.status === 200) {
              // TODO: We will need to change this so that it goes to the MoveTaskOrder view when we're implementing the success UI element in a later story.
              window.location.reload();
            }
          }
        })
        .catch(() => {
          // TODO: Decided if we wnat to display an error notice, log error event, or retry
          setSubmitting(false);
        });
    },
  });

  const handleReviewClick = () => {
    setFilteredShipments(filterShipments(formik.values.shipments));
    setIsModalVisible(true);
  };

  // if showing service items, enable button when shipment and service item are selected
  // if not showing service items, enable button if a shipment is selected
  const isButtonEnabled = moveTaskOrder.availableToPrimeAt
    ? formik.values.shipments.length > 0
    : formik.values.shipments.length > 0 && (formik.values.counselingFee || formik.values.shipmentManagementFee);

  // eslint-disable-next-line camelcase
  const dutyStationPostal = { postal_code: ordersInfo.newDutyStation?.address?.postal_code };

  return (
    <div className={styles.RequestedShipments} data-testid="requested-shipments">
      {shipmentsStatus === 'SUBMITTED' && (
        <>
          <div id="approvalConfirmationModal" style={{ display: isModalVisible ? 'block' : 'none' }}>
            <ShipmentApprovalPreview
              mtoShipments={filteredShipments}
              ordersInfo={ordersInfo}
              allowancesInfo={allowancesInfo}
              customerInfo={customerInfo}
              setIsModalVisible={setIsModalVisible}
              onSubmit={formik.handleSubmit}
              mtoAgents={mtoAgents}
              counselingFee={formik.values.counselingFee}
              shipmentManagementFee={formik.values.shipmentManagementFee}
            />
          </div>

          <form onSubmit={formik.handleSubmit}>
            <h4>Requested shipments</h4>
            <div className={styles.shipmentCards}>
              {mtoShipments &&
                mtoShipments.map((shipment) => (
                  <ShipmentDisplay
                    key={shipment.id}
                    shipmentId={shipment.id}
                    shipmentType={shipment.shipmentType}
                    isSubmitted
                    displayInfo={{
                      heading: shipment.shipmentType === SHIPMENT_OPTIONS.NTS ? 'NTS' : shipment.shipmentType,
                      requestedMoveDate: shipment.requestedPickupDate,
                      currentAddress: shipment.pickupAddress,
                      destinationAddress: shipment.destinationAddress || dutyStationPostal,
                    }}
                    /* eslint-disable-next-line react/jsx-props-no-spreading */
                    {...formik.getFieldProps(`shipments`)}
                  />
                ))}
            </div>

            <div className={styles.serviceItems}>
              {!moveTaskOrder.availableToPrimeAt && (
                <>
                  <h4>Add service items to this move</h4>
                  <Fieldset legend="MTO service items" legendSrOnly id="input-type-fieldset">
                    <Checkbox
                      id="shipmentManagementFee"
                      label="Shipment management fee"
                      name="shipmentManagementFee"
                      onChange={formik.handleChange}
                    />
                    <Checkbox
                      id="counselingFee"
                      label="Counseling fee"
                      name="counselingFee"
                      onChange={formik.handleChange}
                    />
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
                <span>Approve selected shipments</span>
              </Button>
            </div>
          </form>
        </>
      )}

      {shipmentsStatus === 'APPROVED' && (
        <>
          <h4>Approved Shipments</h4>
          <div className={styles.shipmentCards}>
            {mtoShipments &&
              mtoShipments.map((shipment) => (
                <ShipmentDisplay
                  key={shipment.id}
                  shipmentId={shipment.id}
                  shipmentType={shipment.shipmentType}
                  displayInfo={{
                    heading: shipment.shipmentType === SHIPMENT_OPTIONS.NTS ? 'NTS' : shipment.shipmentType,
                    requestedMoveDate: shipment.requestedPickupDate,
                    currentAddress: shipment.pickupAddress,
                    destinationAddress: shipment.destinationAddress || dutyStationPostal,
                  }}
                  isSubmitted={false}
                />
              ))}
          </div>
        </>
      )}

      {shipmentsStatus === 'APPROVED' && (
        <div className={styles.serviceItems}>
          <h4>Service Items</h4>

          <table className="table--stacked">
            <colgroup>
              <col style={{ width: '75%' }} />
              <col style={{ width: '25%' }} />
            </colgroup>
            <tbody>
              {mtoServiceItems &&
                mtoServiceItems
                  .filter((serviceItem) => serviceItem.reServiceCode === 'MS' || serviceItem.reServiceCode === 'CS')
                  .map((serviceItem) => (
                    <tr key={serviceItem.id}>
                      <td data-testid="basicServiceItemName">{serviceItem.reServiceName}</td>
                      <td data-testid="basicServiceItemDate">
                        {serviceItem.status === 'APPROVED' && (
                          <span>
                            <FormCheckmarkIcon className={styles.serviceItemApproval} />{' '}
                            {formatDateFromIso(serviceItem.approvedAt, 'DD MMM YYYY')}
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

RequestedShipments.propTypes = {
  mtoShipments: PropTypes.arrayOf(MTOShipmentShape).isRequired,
  mtoAgents: PropTypes.arrayOf(MTOAgentShape),
  shipmentsStatus: PropTypes.string.isRequired,
  mtoServiceItems: PropTypes.arrayOf(MTOServiceItemShape),
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
      street_address_1: PropTypes.string,
      city: PropTypes.string,
      state: PropTypes.string,
      postal_code: PropTypes.string,
    }),
    backupContactName: PropTypes.string,
    backupContactPhone: PropTypes.string,
    backupContactEmail: PropTypes.string,
  }).isRequired,
  approveMTO: PropTypes.func,
  approveMTOShipment: PropTypes.func,
  moveTaskOrder: MoveTaskOrderShape,
};

RequestedShipments.defaultProps = {
  mtoAgents: [],
  mtoServiceItems: [],
  moveTaskOrder: {},
  approveMTO: () => Promise.resolve(),
  approveMTOShipment: () => Promise.resolve(),
};

export default RequestedShipments;
