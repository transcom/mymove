import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { capitalize } from 'lodash';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { formatDate4DigitYear } from 'shared/formatters';
import PlaceInSitForm from './PlaceInSitForm.jsx';
import './StorageInTransit.css';

export class PlaceInSit extends Component {
  state = {
    showForm: false,
    storageInTransit: {},
  };

  closeForm = () => {
    this.props.onClose();
  };

  componentDidMount() {
    const { estimated_start_date, authorized_start_date } = this.props.sit;
    let startDateValue = estimated_start_date >= authorized_start_date ? estimated_start_date : authorized_start_date;
    this.setState({
      storageInTransit: Object.assign({}, this.props.sit, { actual_start_date: startDateValue }),
    });
  }

  render() {
    const { location, authorized_start_date } = this.props.sit;
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Place into SIT at {capitalize(location)}</div>
        <PlaceInSitForm initialValues={this.state.storageInTransit} />
        <div className="panel-field nested__same-font">
          <div className="usa-input-label unbold">Earliest authorized start</div>
          <div>{formatDate4DigitYear(authorized_start_date)}</div>
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
