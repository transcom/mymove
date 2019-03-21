import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import MoveDocumentView from 'shared/DocumentViewer/MoveDocumentView';
import {
  getAllShipmentDocuments,
  selectShipmentDocument,
  selectShipmentDocuments,
} from 'shared/Entities/modules/shipmentDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { get } from 'lodash';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId, moveDocumentId } = ownProps.match.params;
  const {
    tsp: { shipment: { move = {}, service_member: serviceMember = {} } = {} },
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
    moveDocumentSchema: get(state, 'swaggerPublic.spec.definitions.MoveDocumentPayload', {}),
    moveDocuments: selectShipmentDocuments(state, shipmentId),
    moveLocator: moveLocator || '',
    newDocumentUrl: `/shipments/${shipmentId}/documents/new`,
    serviceMember: { edipi, name },
    uploads: get(shipmentDocument, 'document.uploads', []),
  };
};

const mapDispatchToProps = (dispatch, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  return {
    onDidMount: () => {
      dispatch(loadShipmentDependencies(shipmentId));
      dispatch(getAllShipmentDocuments(shipmentId));
    },
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MoveDocumentView);
