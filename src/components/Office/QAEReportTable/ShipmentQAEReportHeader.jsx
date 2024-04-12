import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { Button } from '@trussworks/react-uswds';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import classnames from 'classnames';
import PropTypes from 'prop-types';

import styles from './ShipmentQAEReportHeader.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatEvaluationReportShipmentAddress } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';
import { milmoveLogger } from 'utils/milmoveLog';
import { createShipmentEvaluationReport } from 'services/ghcApi';
import { SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const ShipmentQAEReportHeader = ({ shipment, destinationDutyLocationPostalCode }) => {
  const { moveCode } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { mutate: createReportMutation } = useMutation(createShipmentEvaluationReport, {
    onSuccess: () => {
      queryClient.invalidateQueries([SHIPMENT_EVALUATION_REPORTS, moveCode]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLogger.error(errorMsg);
    },
  });

  const handleCreateClick = (shipmentID) => {
    createReportMutation(
      { body: { shipmentID }, moveCode },
      {
        onSuccess: (report) => {
          const reportId = report?.id;
          navigate(`/moves/${moveCode}/evaluation-reports/${reportId}`);
        },
      },
    );
  };

  let heading;
  let pickupAddress;
  let destinationAddress;
  let shipmentAccentStyle;

  switch (shipment.shipmentType) {
    case SHIPMENT_OPTIONS.HHG:
      heading = 'HHG';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = shipment?.destinationAddress
        ? formatEvaluationReportShipmentAddress(shipment.destinationAddress)
        : destinationDutyLocationPostalCode;
      shipmentAccentStyle = styles.hhgShipmentType;
      break;
    case SHIPMENT_OPTIONS.NTS:
      heading = 'NTS';
      pickupAddress = formatEvaluationReportShipmentAddress(shipment.pickupAddress);
      destinationAddress = shipment?.storageFacility ? shipment.storageFacility.facilityName : '';
      shipmentAccentStyle = styles.ntsShipmentType;
      break;
    case SHIPMENT_OPTIONS.NTSR:
      heading = 'NTS-Release';
      pickupAddress = shipment?.storageFacility ? shipment.storageFacility.facilityName : '';
      destinationAddress = shipment?.destinationAddress
        ? formatEvaluationReportShipmentAddress(shipment.destinationAddress)
        : destinationDutyLocationPostalCode;
      shipmentAccentStyle = styles.ntsrShipmentType;
      break;
    case SHIPMENT_OPTIONS.PPM:
      heading = 'PPM';
      pickupAddress = shipment.ppmShipment.pickupPostalCode;
      destinationAddress = shipment.ppmShipment.destinationPostalCode;
      shipmentAccentStyle = styles.ppmShipmentType;
      break;
    default:
      break;
  }

  return (
    <>
      <div className={classnames(styles.shipmentAccent, shipmentAccentStyle)} />
      <div className={styles.shipmentInfoContainer}>
        <div data-testid="shipmentHeader" className={styles.shipmentInfo}>
          <h4>
            {heading} Shipment ID: {shipment.shipmentLocator}
          </h4>
          <small>
            {pickupAddress} <FontAwesomeIcon icon="arrow-right" /> {destinationAddress}
          </small>
        </div>
        <Restricted to={permissionTypes.createEvaluationReport}>
          <Button data-testid="shipmentEvaluationCreate" onClick={() => handleCreateClick(shipment.id)}>
            Create report
          </Button>
        </Restricted>
      </div>
    </>
  );
};
ShipmentQAEReportHeader.propTypes = {
  shipment: ShipmentShape.isRequired,
  destinationDutyLocationPostalCode: PropTypes.string.isRequired,
};

export default ShipmentQAEReportHeader;
