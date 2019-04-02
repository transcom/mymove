import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { capitalize } from 'lodash';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faSignInAlt from '@fortawesome/fontawesome-free-solid/faSignInAlt';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { formatDate4DigitYear } from 'shared/formatters';
import PlaceInSitForm from './PlaceInSitForm.jsx';
import './StorageInTransit.css';

export class PlaceInSit extends Component {
  state = {
    showForm: false,
  };

  openForm = () => {
    this.setState({ showForm: true });
  };

  closeForm = () => {
    this.setState({ showForm: false });
  };

  render() {
    const { location, estimated_start_date, authorized_start_date } = this.props.sit;
    const startDatePlaceholder =
      estimated_start_date >= authorized_start_date ? estimated_start_date : authorized_start_date;
    if (this.state.showForm)
      return (
        <div className="storage-in-transit-panel-modal">
          <div className="title">Place into SIT at {capitalize(location)}</div>
          <PlaceInSitForm />
          <div className="panel-field">
            <span className="field-title unbold">Earliest authorized start</span>
            <span>{formatDate4DigitYear(startDatePlaceholder)}</span>
          </div>
          <div className="usa-grid-full align-center-vertical">
            <div className="usa-width-one-half">
              <p className="cancel-link">
                <a className="usa-button-secondary" onClick={this.closeForm}>
                  Cancel
                </a>
              </p>
            </div>
            <div className="usa-width-one-half align-right">
              <button
                className="button usa-button-primary storage-in-transit-request-form-send-request-button"
                disabled={!this.props.formEnabled}
              >
                Approve
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <span className="approve-sit">
        <a className="approve-sit-link" onClick={this.openForm}>
          <FontAwesomeIcon className="icon" icon={faSignInAlt} />
          Place into SIT
        </a>
      </span>
    );
  }
}

PlaceInSit.propTypes = { sit: PropTypes.object.isRequired };

function mapStateToProps(state) {
  return {};
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PlaceInSit);
