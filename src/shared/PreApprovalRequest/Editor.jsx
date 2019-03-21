import { cloneDeep } from 'lodash';
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import PreApprovalForm, { formName as PreApprovalFormName } from 'shared/PreApprovalRequest/PreApprovalForm.jsx';
import {
  formatToBaseQuantity,
  convertFromBaseQuantity,
  formatToDimensionsInches,
  formatDimensionsToThousandthInches,
  formatCents,
} from 'shared/formatters';
import { submit, isValid, isDirty, isSubmitting, hasSubmitSucceeded } from 'redux-form';
import { get } from 'lodash';
import { convertDollarsToCents } from 'shared/utils';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
export class Editor extends Component {
  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) this.props.onSaveComplete();
  }
  onSubmit = values => {
    // Convert quantity_1 to base quantity unit before hitting endpoint
    if (values.quantity_1) {
      values.quantity_1 = formatToBaseQuantity(values.quantity_1);
    }
    values.tariff400ng_item_id = values.tariff400ng_item.id;

    values.estimate_amount_cents = convertDollarsToCents(get(values, 'estimate_amount_cents'));
    values.actual_amount_cents = convertDollarsToCents(get(values, 'actual_amount_cents'));

    formatDimensionsToThousandthInches(values.item_dimensions);
    formatDimensionsToThousandthInches(values.crate_dimensions);

    this.props.saveEdit(this.props.shipmentLineItem.id, values);
  };
  render() {
    let initialValues = cloneDeep(this.props.shipmentLineItem);

    // Leaving quantity_2 alone in the swagger definition until we know for sure we won't use
    // Delete quantities so we don't unnecessarily send them back
    delete initialValues.quantity_2;
    if (!initialValues.quantity_1 || initialValues.quantity_1 <= 0) {
      delete initialValues.quantity_1;
    } else {
      initialValues.quantity_1 = convertFromBaseQuantity(initialValues.quantity_1);
    }

    initialValues.estimate_amount_cents = formatCents(initialValues.estimate_amount_cents);
    initialValues.actual_amount_cents = formatCents(initialValues.actual_amount_cents);
    initialValues.item_dimensions = formatToDimensionsInches(initialValues.item_dimensions);
    initialValues.crate_dimensions = formatToDimensionsInches(initialValues.crate_dimensions);

    return (
      <div className="pre-approval-panel-modal pre-approval-edit">
        <div className="title">Edit request</div>
        <PreApprovalForm
          tariff400ngItems={this.props.tariff400ngItems}
          onSubmit={this.onSubmit}
          initialValues={initialValues}
        />
        <div className="usa-grid-full align-center-vertical">
          <div className="usa-width-one-half">
            <p className="cancel-link">
              <a className="usa-button-secondary" onClick={this.props.cancelEdit}>
                Cancel
              </a>
            </p>
          </div>

          <div className="usa-width-one-half align-right">
            <button
              className="button usa-button-primary"
              disabled={!this.props.formEnabled}
              onClick={this.props.submitForm}
            >
              Save
            </button>
          </div>
        </div>
      </div>
    );
  }
}
Editor.propTypes = {
  tariff400ngItems: PropTypes.array,
  shipmentLineItem: PropTypes.object.isRequired,
  saveEdit: PropTypes.func.isRequired,
  cancelEdit: PropTypes.func.isRequired,
  onSaveComplete: PropTypes.func.isRequired,
  formEnabled: PropTypes.bool.isRequired,
  hasSubmitSucceeded: PropTypes.bool.isRequired,
  submitForm: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled:
      isDirty(PreApprovalFormName)(state) &&
      isValid(PreApprovalFormName)(state) &&
      !isSubmitting(PreApprovalFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(PreApprovalFormName)(state),
  };
}

function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators(
    {
      submitForm: () => submit(PreApprovalFormName),
    },
    dispatch,
  );
}
export default connect(mapStateToProps, mapDispatchToProps)(Editor);
