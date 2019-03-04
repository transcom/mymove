import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

import BasicPanel from 'shared/BasicPanel';
import Alert from 'shared/Alert';
import Creator from 'shared/StorageInTransit/Creator';
import { selectStorageInTransits, createStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import { loadEntitlements } from '../../scenes/TransportationServiceProvider/ducks';

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

  onFormActivation = isFormActive => {
    this.setState({ isRequestActionable: !isFormActive });
  };

  onSubmit = createPayload => {
    return this.props.createStorageInTransit(this.props.shipmentId, createPayload);
  };

  render() {
    const { storageInTransitEntitlement, storageInTransits } = this.props;
    const { error, isCreatorActionable } = this.state;
    const daysUsed = 0; // placeholder
    const daysRemaining = storageInTransitEntitlement - daysUsed;
    return (
      <div className="storage-in-transit-panel">
        <BasicPanel title="Storage in Transit (SIT)">
          {error && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <div className="column-head">
            Entitlement: {storageInTransitEntitlement} days <span className="unbold">({daysRemaining} remaining)</span>
          </div>
          {storageInTransits !== undefined &&
            storageInTransits.map(row => {
              return (
                <div key={row.id} className="column-head">
                  {row.location.charAt(0) + row.location.slice(1).toLowerCase()} SIT
                  <span className="unbold">
                    {' '}
                    <span id={'sit-status-text'}>Status:</span>{' '}
                    <FontAwesomeIcon className="icon icon-grey" icon={faClock} />
                  </span>
                  <span>SIT {row.status.charAt(0) + row.status.slice(1).toLowerCase()} </span>
                </div>
              );
            })}
          {isCreatorActionable && (
            <Creator onFormActivation={this.onFormActivation} saveStorageInTransit={this.onSubmit} />
          )}
        </BasicPanel>
      </div>
    );
  }
}

StorageInTransitPanel.propTypes = {
  storageInTransits: PropTypes.array,
  shipmentId: PropTypes.string,
  storageInTransitEntitlement: PropTypes.number,
};

function mapStateToProps(state, ownProps) {
  return {
    storageInTransits: selectStorageInTransits(state, ownProps.shipmentId),
    storageInTransitEntitlement: loadEntitlements(state).storage_in_transit,
    shipmentId: ownProps.shipmentId,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createStorageInTransit }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(StorageInTransitPanel);
