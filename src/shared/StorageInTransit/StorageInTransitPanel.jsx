import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';

import BasicPanel from 'shared/BasicPanel';
import Alert from 'shared/Alert';
import Creator from 'shared/StorageInTransit/Creator';

export class StorageInTransitPanel extends Component {
  constructor() {
    super();
    this.state = {
      isRequestActionable: true,
      isCreatorActionable: true,
      error: null,
    };
  }
  closeError = () => {
    this.setState({ error: null });
  };

  render() {
    return (
      <div className="storage-in-transit-panel">
        <BasicPanel title={'Storage in Transit (SIT)'}>
          {this.state.error && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <div className="column-subhead">Entitlement: {this.props.SITEntitlement} days</div>
          {this.state.isCreatorActionable && <Creator SITRequests={this.props.SITRequests} />}
        </BasicPanel>
      </div>
    );
  }
}

StorageInTransitPanel.propTypes = {
  SITRequests: PropTypes.array,
  shipmentId: PropTypes.string,
  SITEntitlement: PropTypes.number,
};

/*
function mapStateToProps(state, ownProps) {
  return {
    SITRequests: selectSortedSITRequests(state, ownProps.shipmentId),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { createSITRequest, deleteSITRequest, approveSITRequest, updateSITRequest },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(StorageInTransitPanel);
*/
export default connect(null, null)(StorageInTransitPanel);
