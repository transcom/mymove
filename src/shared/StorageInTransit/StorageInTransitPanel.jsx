import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { some } from 'lodash';

import BasicPanel from 'shared/BasicPanel';
import Alert from 'shared/Alert';
import StorageInTransit from 'shared/StorageInTransit/StorageInTransit';
import Creator from 'shared/StorageInTransit/Creator';
import { selectStorageInTransits, createStorageInTransit } from 'shared/Entities/modules/storageInTransits';
import { calculateEntitlementsForShipment } from 'shared/Entities/modules/shipments';
import { calculateEntitlementsForMove } from 'shared/Entities/modules/moves';

import { isTspSite } from 'shared/constants.js';
import SitStatusIcon from './SitStatusIcon';
import { sitTotalDaysUsed } from 'shared/StorageInTransit/calculator';

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
    const hasRequestedSIT = some(storageInTransits, sit => sit.status === 'REQUESTED');
    const hasInSIT = some(storageInTransits, sit => sit.status === 'IN_SIT');
    const daysRemaining = storageInTransitEntitlement - sitTotalDaysUsed(storageInTransits);

    return (
      <div className="storage-in-transit-panel" data-cy="storage-in-transit-panel">
        <BasicPanel
          title="Storage in Transit (SIT)"
          titleExtension={!isTspSite && hasRequestedSIT ? <SitStatusIcon isTspSite={isTspSite} /> : null}
        >
          {error && (
            <Alert type="error" heading="Oops, something went wrong!" onRemove={this.closeError}>
              <span className="warning--header">Please refresh the page and try again.</span>
            </Alert>
          )}
          <div className="column-head">
            Entitlement: {storageInTransitEntitlement} days
            {hasInSIT && <span className="unbold"> ({daysRemaining} remaining)</span>}
          </div>
          {storageInTransits !== undefined &&
            storageInTransits.map(storageInTransit => {
              return (
                <StorageInTransit
                  key={storageInTransit.id}
                  storageInTransit={storageInTransit}
                  daysRemaining={daysRemaining}
                />
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
function getStorageInTransitEntitlement(state, resourceId) {
  let storageInTransitEntitlement = 0;
  if (isTspSite) {
    storageInTransitEntitlement = calculateEntitlementsForShipment(state, resourceId).storage_in_transit;
  } else {
    storageInTransitEntitlement = calculateEntitlementsForMove(state, resourceId).storage_in_transit;
  }
  return storageInTransitEntitlement;
}

function mapStateToProps(state, ownProps) {
  const resourceId = isTspSite ? ownProps.shipmentId : ownProps.moveId;

  return {
    storageInTransits: selectStorageInTransits(state, ownProps.shipmentId),
    storageInTransitEntitlement: getStorageInTransitEntitlement(state, resourceId),
    shipmentId: ownProps.shipmentId,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ createStorageInTransit }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(StorageInTransitPanel);
