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
import { formatDate4DigitYear } from 'shared/formatters';
import { calculateEntitlementsForMove } from 'shared/Entities/modules/moves';

import { isTspSite } from 'shared/constants.js';

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
            storageInTransits.map(storageInTransit => {
              return (
                <div key={storageInTransit.id} className="storage-in-transit">
                  <div className="column-head">
                    {storageInTransit.location.charAt(0) + storageInTransit.location.slice(1).toLowerCase()} SIT
                    <span className="unbold">
                      {' '}
                      <span id="sit-status-text">Status:</span>{' '}
                      <FontAwesomeIcon className="icon icon-grey" icon={faClock} />
                    </span>
                    <span>
                      SIT {storageInTransit.status.charAt(0) + storageInTransit.status.slice(1).toLowerCase()}{' '}
                    </span>
                  </div>
                  <div className="usa-width-one-whole">
                    <div className="usa-width-one-half">
                      <div className="column-subhead">Dates</div>
                      <div className="panel-field">
                        <span className="field-title unbold">Est. start date</span>
                        <span className="field-value">
                          {formatDate4DigitYear(storageInTransit.estimated_start_date)}
                        </span>
                      </div>
                      {storageInTransit.notes !== undefined && (
                        <div className="sit-notes">
                          <div className="column-subhead">Note</div>
                          <div className="panel-field">
                            <span className="field-title unbold">{storageInTransit.notes}</span>
                          </div>
                        </div>
                      )}
                    </div>
                    <div className="usa-width-one-half">
                      <div className="column-subhead">Warehouse</div>
                      <div className="panel-field">
                        <span className="field-title unbold">Warehouse ID</span>
                        <span className="field-value">{storageInTransit.warehouse_id}</span>
                      </div>
                      <div className="panel-field">
                        <span className="field-title unbold">Contact info</span>
                        <span className="field-value">
                          {storageInTransit.warehouse_name}
                          <br />
                          {storageInTransit.warehouse_address.street_address_1}
                          <br />
                          {storageInTransit.warehouse_address.city}, {storageInTransit.warehouse_address.state}{' '}
                          {storageInTransit.warehouse_address.postal_code}
                          <br />
                          {storageInTransit.warehouse_phone}
                          <br />
                          {storageInTransit.warehouse_email}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          {isCreatorActionable &&
            isTspSite && <Creator onFormActivation={this.onFormActivation} saveStorageInTransit={this.onSubmit} />}
        </BasicPanel>
      </div>
    );
  }
}

StorageInTransitPanel.propTypes = {
  storageInTransits: PropTypes.array,
  shipmentId: PropTypes.string,
  storageInTransitEntitlement: PropTypes.number,
  moveId: PropTypes.string,
};

// Service Member entitlement is stored different depending on if this is
// being called from the TSP or Office site. Need to check for both.
// moveId is needed to find the entitlement on the Office side. It is
// not needed to pull entitlement from the TSP side.
// calculateEntitlementsForMove is a more up-to-date way of storing data
function getStorageInTransitEntitlement(state, moveId) {
  let storageInTransitEntitlement = 0;
  if (isTspSite) {
    storageInTransitEntitlement = loadEntitlements(state).storage_in_transit;
  } else {
    storageInTransitEntitlement = calculateEntitlementsForMove(state, moveId).storage_in_transit;
  }
  return storageInTransitEntitlement;
}

function mapStateToProps(state, ownProps) {
  const moveId = ownProps.moveId;
  return {
    storageInTransits: selectStorageInTransits(state, ownProps.shipmentId),
    storageInTransitEntitlement: getStorageInTransitEntitlement(state, moveId),
    shipmentId: ownProps.shipmentId,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createStorageInTransit }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(StorageInTransitPanel);
