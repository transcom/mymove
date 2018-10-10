import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import NewDocumentView from 'shared/DocumentViewer/NewDocumentView';
import {
  getAllShipmentDocuments,
  getShipmentDocumentsLabel,
  createShipmentDocumentLabel,
  createShipmentDocument,
} from 'shared/Entities/modules/shipmentDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { get } from 'lodash';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  const {
    tsp: { shipment: { move = {}, service_member: serviceMember = {} } = {} },
    entities: { moveDocuments = {}, uploads = {} },
  } = state;
  const { locator: moveLocator } = move;
  const { edipi = '' } = serviceMember;
  const name = stringifyName(serviceMember);

  return {
    // moveDocumentSchema: get(
    // state,
    // 'swagger.spec.definitions.MoveDocumentPayload',
    // {},
    // ),
    shipmentId,
    genericMoveDocSchema: get(
      state,
      'swaggerPublic.spec.definitions.CreateGenericMoveDocumentPayload',
      {},
    ),
    moveDocSchema: get(
      state,
      'swaggerPublic.spec.definitions.MoveDocumentPayload',
      {},
    ),
    moveDocuments: Object.values(moveDocuments),
    moveLocator: moveLocator || '',
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
    createShipmentDocument: (shipmentId, body) =>
      dispatch(
        createShipmentDocument(createShipmentDocumentLabel, shipmentId, body),
      ),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(NewDocumentView);
