import React from 'react';
import * as PropTypes from 'prop-types';

import { ShipmentShape } from 'types/shipment';
import ShipmentInfoList from 'components/Office/DefinitionLists/ShipmentInfoList';
import PPMShipmentInfoList from 'components/Office/DefinitionLists/PPMShipmentInfoList';
import NTSRShipmentInfoList from 'components/Office/DefinitionLists/NTSRShipmentInfoList';
import NTSShipmentInfoList from 'components/Office/DefinitionLists/NTSShipmentInfoList';
import { SHIPMENT_OPTIONS } from 'shared/constants';

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
  warnIfMissing: PropTypes.arrayOf(PropTypes.string),
  errorIfMissing: PropTypes.arrayOf(PropTypes.string),
  showWhenCollapsed: PropTypes.arrayOf(PropTypes.string),
  neverShow: PropTypes.arrayOf(PropTypes.string),
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
    SHIPMENT_OPTIONS.NTSR,
    SHIPMENT_OPTIONS.PPM,
  ]),
  isForEvaluationReport: PropTypes.bool,
  destinationDutyLocationPostalCode: PropTypes.string,
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
};

export default ShipmentInfoListSelector;
