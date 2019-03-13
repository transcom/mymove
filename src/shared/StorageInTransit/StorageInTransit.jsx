import React, { Component } from 'react';
import PropTypes from 'prop-types';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

import './StorageInTransit.css';
import { formatDate4DigitYear } from 'shared/formatters';

export class StorageInTransit extends Component {
  render() {
    const { storageInTransit } = this.props;

    return (
      <div key={storageInTransit.id} className="storage-in-transit">
        <div className="column-head">
          {storageInTransit.location.charAt(0) + storageInTransit.location.slice(1).toLowerCase()} SIT
          <span className="unbold">
            {' '}
            <span id={'sit-status-text'}>Status:</span> <FontAwesomeIcon className="icon icon-grey" icon={faClock} />
          </span>
          <span>SIT {storageInTransit.status.charAt(0) + storageInTransit.status.slice(1).toLowerCase()} </span>
        </div>
        <div className="usa-width-one-whole">
          <div className="usa-width-one-half">
            <div className="column-subhead">Dates</div>
            <div className="panel-field">
              <span className="field-title unbold">Est. start date</span>
              <span className="field-value">{formatDate4DigitYear(storageInTransit.estimated_start_date)}</span>
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
  }
}

StorageInTransit.propTypes = {
  storageInTransit: PropTypes.object.isRequired,
};

export default StorageInTransit;
