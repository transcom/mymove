import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './StorageInTransit.css';

import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitForm, { formName as StorageInTransitFormName } from './StorageInTransitForm.jsx';
import StorageInTransitTspEditForm, {
  formName as StorageInTransitTspEditFormName,
} from './StorageInTransitTspEditForm.jsx';
import { getPublicShipment } from 'shared/Entities/modules/shipments';
import { updateStorageInTransit } from 'shared/Entities/modules/storageInTransits';

export class TspEditor extends Component {
  state = {
    closeOnSubmit: true,
  };

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.props.onClose();
    }
  }

  closeForm = () => {
    this.props.onClose();
  };

  saveAndClose = () => {
    this.setState({ closeOnSubmit: true }, () => {
      this.props.submitForm();
    });
  };

  onSubmit = values => {
    const { storageInTransit } = this.props;
    this.props.updateStorageInTransit(storageInTransit.shipment_id, storageInTransit.id, values).then(() => {
      this.props.getPublicShipment(storageInTransit.shipment_id);
    });
  };

  render() {
    const isRequested = this.props.storageInTransit.status === 'REQUESTED';

    return (
      <div className="storage-in-transit-panel-modal">
        <div className="editable-panel is-editable">
          <div className="title">Edit SIT Request</div>
          {isRequested ? (
            <StorageInTransitForm onSubmit={this.onSubmit} initialValues={this.props.storageInTransit} />
          ) : (
            <StorageInTransitTspEditForm
              minDate={this.props.storageInTransit.authorized_start_date}
              onSubmit={this.onSubmit}
              initialValues={this.props.storageInTransit}
            />
          )}
          <div className="usa-grid-full">
            <div className="usa-width-one-half">
              <p className="cancel-link">
                <a className="usa-button-secondary" onClick={this.closeForm}>
                  Cancel
                </a>
              </p>
            </div>
            <div className="usa-width-one-half align-right">
              <button
                className="button usa-button-primary"
                disabled={!this.props.formEnabled}
                onClick={this.saveAndClose}
              >
                Save
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

TspEditor.propTypes = {
  updateStorageInTransit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
  storageInTransit: PropTypes.object.isRequired,
  submitForm: PropTypes.func.isRequired,
};

function formNameSelector(props) {
  const storageInTransit = props.storageInTransit;
  let formName = '';
  storageInTransit.status === 'REQUESTED'
    ? (formName = StorageInTransitFormName)
    : (formName = StorageInTransitTspEditFormName);
  return formName;
}

function mapStateToProps(state, props) {
  return {
    formEnabled: isValid(formNameSelector(props))(state) && !isSubmitting(formNameSelector(props))(state),
    hasSubmitSucceeded: hasSubmitSucceeded(formNameSelector(props))(state),
  };
}

function mapDispatchToProps(dispatch, ownProps) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(formNameSelector(ownProps)),
      updateStorageInTransit,
      getPublicShipment,
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(TspEditor);
