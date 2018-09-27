import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import MoveDocumentView from 'shared/DocumentViewer/MoveDocumentView';
import {
  getAllShipmentDocuments,
  getShipmentDocumentsLabel,
} from 'shared/Entities/modules/shipmentDocuments';

const mapStateToProps = state => {
  const {
    tsp: { shipment = {} },
    entities: { moveDocuments = [] },
  } = state;
  return { shipment, moveDocuments };
};

const mapDispatchToProps = dispatch => ({
  onDidMount: shipmentId => {
    dispatch(loadShipmentDependencies(shipmentId));
    dispatch(getAllShipmentDocuments(getShipmentDocumentsLabel, shipmentId));
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(MoveDocumentView);
