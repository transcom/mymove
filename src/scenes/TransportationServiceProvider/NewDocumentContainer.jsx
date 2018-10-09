import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import NewDocumentView from 'shared/DocumentViewer/NewDocumentView';
import {
  getAllShipmentDocuments,
  getShipmentDocumentsLabel,
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

  return {
    // documentDetailUrlPrefix: `/shipments/${shipmentId}/documents`,
    // moveDocumentSchema: get(
    // state,
    // 'swagger.spec.definitions.MoveDocumentPayload',
    // {},
    // ),
    moveDocuments: Object.values(moveDocuments),
    moveLocator: moveLocator || '',
    // newDocumentUrl: `/shipments/${shipmentId}/documents/new`,
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

export default connect(mapStateToProps, mapDispatchToProps)(NewDocumentView);
