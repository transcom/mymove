import React from 'react';
import * as PropTypes from 'prop-types';

import { ShipmentShape } from 'types/shipment';
import ShipmentInfoList from 'components/Office/DefinitionLists/ShipmentInfoList';
import PPMShipmentInfoList from 'components/Office/DefinitionLists/PPMShipmentInfoList';
import NTSRShipmentInfoList from 'components/Office/DefinitionLists/NTSRShipmentInfoList';
import NTSShipmentInfoList from 'components/Office/DefinitionLists/NTSShipmentInfoList';
import BoatShipmentInfoList from 'components/Office/DefinitionLists/BoatShipmentInfoList';
import { SHIPMENT_OPTIONS, SHIPMENT_TYPES } from 'shared/constants';
import { fieldValidationShape } from 'utils/displayFlags';

const ShipmentInfoListSelector = ({
  className,
  shipment,
  isExpanded,
  warnIfMissing,
  errorIfMissing,
  showWhenCollapsed,
  neverShow,
  shipmentType,
  isForEvaluationReport,
  destinationDutyLocationPostalCode,
  onErrorModalToggle,
}) => {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.PPM:
      return (
        <PPMShipmentInfoList
          className={className}
          shipment={shipment}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          shipmentType={shipmentType}
          showWhenCollapsed={showWhenCollapsed}
          isExpanded={isExpanded}
          isForEvaluationReport={isForEvaluationReport}
          onErrorModalToggle={onErrorModalToggle}
        />
      );
    case SHIPMENT_OPTIONS.HHG:
      return (
        <ShipmentInfoList
          className={className}
          shipment={shipment}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          shipmentType={shipmentType}
          showWhenCollapsed={showWhenCollapsed}
          isForEvaluationReport={isForEvaluationReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      );
    case SHIPMENT_OPTIONS.NTSR:
      return (
        <NTSRShipmentInfoList
          className={className}
          shipment={shipment}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
          neverShow={neverShow}
          isForEvaluationReport={isForEvaluationReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      );
    case SHIPMENT_OPTIONS.NTS:
      return (
        <NTSShipmentInfoList
          className={className}
          shipment={shipment}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          showWhenCollapsed={showWhenCollapsed}
          neverShow={neverShow}
          isForEvaluationReport={isForEvaluationReport}
        />
      );
    case SHIPMENT_OPTIONS.MOBILE_HOME:
      return (
        <MobileHomeShipmentInfoList
          className={className}
          shipment={shipment}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          shipmentType={shipmentType}
          showWhenCollapsed={showWhenCollapsed}
          isForEvaluationReport={isForEvaluationReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      );
    case SHIPMENT_OPTIONS.BOAT:
    case SHIPMENT_TYPES.BOAT_HAUL_AWAY:
    case SHIPMENT_TYPES.BOAT_TOW_AWAY:
      return (
        <BoatShipmentInfoList
          className={className}
          shipment={shipment}
          isExpanded={isExpanded}
          warnIfMissing={warnIfMissing}
          errorIfMissing={errorIfMissing}
          shipmentType={shipmentType}
          showWhenCollapsed={showWhenCollapsed}
          isForEvaluationReport={isForEvaluationReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      );
    default:
      return (
        <ShipmentInfoList
          className={className}
          shipment={shipment}
          shipmentType={shipmentType}
          isExpanded={isExpanded}
          isForEvaluationReport={isForEvaluationReport}
          destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
        />
      );
  }
};

ShipmentInfoListSelector.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
  warnIfMissing: PropTypes.arrayOf(fieldValidationShape),
  errorIfMissing: PropTypes.arrayOf(fieldValidationShape),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_TYPES.BOAT_HAUL_AWAY,
    SHIPMENT_TYPES.BOAT_TOW_AWAY,
    SHIPMENT_OPTIONS.MOBILE_HOME,
  ]),
  isForEvaluationReport: PropTypes.bool,
  destinationDutyLocationPostalCode: PropTypes.string,
  onErrorModalToggle: PropTypes.func,
};

ShipmentInfoListSelector.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
  className: '',
  isExpanded: false,
  warnIfMissing: [],
  errorIfMissing: [],
  showWhenCollapsed: [],
  neverShow: [],
  isForEvaluationReport: false,
  destinationDutyLocationPostalCode: '',
  onErrorModalToggle: undefined,
};

export default ShipmentInfoListSelector;
