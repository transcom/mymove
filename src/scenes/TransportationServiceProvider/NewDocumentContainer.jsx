import { connect } from 'react-redux';
import NewDocumentView from 'shared/DocumentViewer/NewDocumentView';
import {
  getAllShipmentDocuments,
  createShipmentDocument,
  selectShipmentDocuments,
} from 'shared/Entities/modules/shipmentDocuments';
import { selectShipment, getPublicShipment } from 'shared/Entities/modules/shipments';
import { stringifyName } from 'shared/utils/serviceMember';
import { get } from 'lodash';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  const {
    entities: { uploads = {} },
  } = state;
  const shipment = selectShipment(state, shipmentId);
  const serviceMember = shipment.service_member || {};
  const { locator: moveLocator } = shipment.move || {};
  const { edipi = '' } = serviceMember;
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
      dispatch(getPublicShipment(shipmentId));
      dispatch(getAllShipmentDocuments(shipmentId));
    },
    createShipmentDocument: (shipmentId, body) => dispatch(createShipmentDocument(shipmentId, body)),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(NewDocumentView);
