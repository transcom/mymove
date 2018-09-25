import { connect } from 'react-redux';
import { loadShipmentDependencies } from './ducks';
import MoveDocumentView from 'shared/DocumentViewer/MoveDocumentView';

const mapStateToProps = state => {
  const {
    tsp: { shipment = {} },
  } = state;
  return { shipment: shipment };
};

const mapDispatchToProps = dispatch => ({
  onDidMount: shipmentId => {
    dispatch(loadShipmentDependencies(shipmentId));
  },
});

export default connect(mapStateToProps, mapDispatchToProps)(MoveDocumentView);
