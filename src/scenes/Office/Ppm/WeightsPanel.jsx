import React from 'react';
import { connect } from 'react-redux';
import {
  selectAllDocumentsForMove,
  findOKedVehicleWeightTickets,
  findOKedProgearWeightTickets,
  findPendingWeightTickets,
} from 'shared/Entities/modules/moveDocuments';
import Alert from 'shared/Alert';

import { PanelField } from 'shared/EditablePanel';
import { editablePanelify } from 'shared/EditablePanel';
import { formatWeight } from 'utils/formatters';

function sumWeights(moveDocs) {
  return moveDocs.reduce(function (sum, { empty_weight, full_weight }) {
    // empty_weight and full_weight can be blank
    empty_weight = empty_weight || 0;
    full_weight = full_weight || 0;

    // Minimize the damage from having an empty_weight that is larger than the full_weight.
    if (empty_weight > full_weight) {
      return 0;
    }

    return sum + full_weight - empty_weight;
  }, 0);
}

const WeightDisplay = ({
  hasPendingWeightTickets,
  ppmPaymentRequestedFlag,
  vehicleWeightTicketWeight,
  progearWeightTicketWeight,
}) => {
  return (
    <>
      {ppmPaymentRequestedFlag && hasPendingWeightTickets && (
        <div className="missing-info-alert">
          <Alert type="warning">There are more weight tickets awaiting review.</Alert>
        </div>
      )}
      <div className="editable-panel-column">
        <PanelField title="Net Weight" value={formatWeight(vehicleWeightTicketWeight)} />
      </div>
      <div className="editable-panel-column">
        <PanelField title="Pro-Gear" value={formatWeight(progearWeightTicketWeight)} />
      </div>
    </>
  );
};

const WeightPanel = editablePanelify(WeightDisplay, null, false);

function mapStateToProps(state, ownProps) {
  const moveDocs = selectAllDocumentsForMove(state, ownProps.moveId);

  return {
    ppmPaymentRequestedFlag: true,
    vehicleWeightTicketWeight: sumWeights(findOKedVehicleWeightTickets(moveDocs)),
    progearWeightTicketWeight: sumWeights(findOKedProgearWeightTickets(moveDocs)),
    hasPendingWeightTickets: findPendingWeightTickets(moveDocs).length > 0,
  };
}

export default connect(mapStateToProps)(WeightPanel);
