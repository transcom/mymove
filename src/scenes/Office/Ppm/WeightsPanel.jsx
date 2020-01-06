import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import { selectAllDocumentsForMove, findPendingWeightTickets } from 'shared/Entities/modules/moveDocuments';
import Alert from 'shared/Alert';

import { PanelSwaggerField } from 'shared/EditablePanel';
import { editablePanelify } from 'shared/EditablePanel';

const WeightDisplay = ({ ppmSchema, ppm, hasWeightTicketsPending, ppmPaymentRequestedFlag }) => {
  const fieldProps = {
    schema: ppmSchema,
    values: ppm,
  };
  return (
    <div className="editable-panel-column">
      {ppmPaymentRequestedFlag && hasWeightTicketsPending && (
        <div className="missing-info-alert">
          <Alert type="warning">There are more weight tickets awaiting review.</Alert>
        </div>
      )}
      <PanelSwaggerField title="Net Weight" fieldName="net_weight" required {...fieldProps} />
    </div>
  );
};

const WeightPanel = editablePanelify(WeightDisplay, null, false);

function mapStateToProps(state, ownProps) {
  const ppm = selectPPMForMove(state, ownProps.moveId);
  const moveDocs = selectAllDocumentsForMove(state, ownProps.moveId);
  const hasWeightTicketsPending = findPendingWeightTickets(moveDocs).length > 0;

  return {
    // Wrapper
    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,
    hasWeightTicketsPending,
  };
}

export default connect(mapStateToProps)(WeightPanel);
