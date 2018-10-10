import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import MoveDocumentView from 'shared/DocumentViewer/MoveDocumentView';
import {
  getAllShipmentDocuments,
  getShipmentDocumentsLabel,
  selectShipmentDocument,
} from 'shared/Entities/modules/shipmentDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { get } from 'lodash';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId, moveDocumentId } = ownProps.match.params;
  const {
    tsp: { shipment: { move = {}, service_member: serviceMember = {} } = {} },
    entities: { moveDocuments = {}, uploads = {} },
  } = state;
  const { locator: moveLocator } = move;
  const { edipi = '' } = serviceMember;
  const name = stringifyName(serviceMember);
  const shipmentDocument = selectShipmentDocument(state, moveDocumentId);

  return {
    documentDetailUrlPrefix: `/shipments/${shipmentId}/documents`,
    moveDocument: {
      createdAt: get(shipmentDocument, 'document.uploads.0.created_at', ''),
      type: shipmentDocument.move_document_type,
      notes: shipmentDocument.notes || '',
      ...shipmentDocument,
    },
    moveDocumentSchema: get(
      state,
      'swaggerPublic.spec.definitions.MoveDocumentPayload',
      {},
    ),
    moveDocuments: Object.values(moveDocuments),
    moveLocator: moveLocator || '',
    newDocumentUrl: `/shipments/${shipmentId}/documents/new`,
    serviceMember: { edipi, name },
    uploads: Object.values(uploads),
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  return {
    onDidMount: () => {
      dispatch(loadShipmentDependencies(shipmentId));
      dispatch(getAllShipmentDocuments(getShipmentDocumentsLabel, shipmentId));
    },
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MoveDocumentView);
