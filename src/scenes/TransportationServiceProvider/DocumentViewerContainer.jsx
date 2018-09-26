import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import MoveDocumentView from 'shared/DocumentViewer/MoveDocumentView';
import {
  getAllShipmentDocuments,
  getShipmentDocumentsLabel,
} from 'shared/Entities/modules/shipmentDocuments';

const mapStateToProps = (state, ownProps) => {
  const { shipmentId } = ownProps.match.params;
  const {
    tsp: { shipment: { move = {}, service_member: serviceMember = {} } = {} },
    entities: { moveDocuments = {}, uploads = {} },
  } = state;
  const { locator: moveLocator } = move;
  const {
    edipi = '',
    last_name: lastName = '',
    first_name: firstName = '',
  } = serviceMember;
  const name = [lastName, firstName].filter(name => !!name).join(', ');

  return {
    documentDetailUrlPrefix: `/shipments/${shipmentId}/documents`,
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
