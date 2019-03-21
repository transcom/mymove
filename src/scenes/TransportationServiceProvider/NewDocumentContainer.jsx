import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import NewDocumentView from 'shared/DocumentViewer/NewDocumentView';
import {
  getAllShipmentDocuments,
  createShipmentDocument,
  selectShipmentDocuments,
} from 'shared/Entities/modules/shipmentDocuments';
import { stringifyName } from 'shared/utils/serviceMember';
import { get } from 'lodash';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  const {
    tsp,
    entities: { uploads = {} },
  } = state;
  const serviceMember = get(tsp, 'serviceMember', {});
  const { locator: moveLocator } = get(tsp, 'shipment.move', {});
  const { edipi = '' } = get(tsp, 'serviceMember', {});
  const name = stringifyName(serviceMember);

  return {
    shipmentId,
    genericMoveDocSchema: get(state, 'swaggerPublic.spec.definitions.CreateGenericMoveDocumentPayload', {}),
    moveDocSchema: get(state, 'swaggerPublic.spec.definitions.MoveDocumentPayload', {}),
    moveDocuments: selectShipmentDocuments(state, shipmentId),
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
      dispatch(getAllShipmentDocuments(shipmentId));
    },
    createShipmentDocument: (shipmentId, body) => dispatch(createShipmentDocument(shipmentId, body)),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(NewDocumentView);
