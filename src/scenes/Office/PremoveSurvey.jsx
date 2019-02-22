// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import classNames from 'classnames';
import { reduxForm, FormSection, getFormValues } from 'redux-form';

import { PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

import './office.css';
import { getRequestStatus } from 'shared/Swagger/selectors';
import { humanReadableError } from 'shared/utils';
import Alert from 'shared/Alert';

const surveyFields = [
  'pm_survey_conducted_date',
  'pm_survey_planned_pack_date',
  'pm_survey_planned_pickup_date',
  'pm_survey_planned_delivery_date',
  'pm_survey_weight_estimate',
  'pm_survey_progear_weight_estimate',
  'pm_survey_spouse_progear_weight_estimate',
  'pm_survey_notes',
  'pm_survey_method',
];

const premoveSurveyUpdateShipmentLabel = 'shipment.updateShipment.premoveSurvey';

// TODO: Refactor when we switch to using a wizard
// Editable panel specific to Enter pre-move survey. Due to not using a wizard to enter the pre-move survey this
// panel has highly specific behavior (opening the edit view via clicking on Enter pre-move survey button)
export class PreMoveSurveyEditablePanel extends Component {
  handleEditClick = e => {
    e.preventDefault();
    this.props.onEdit(true);
  };
  handleCancelClick = e => {
    e.preventDefault();
    this.props.onCancel();
  };
  handleSaveClick = e => {
    e.preventDefault();
    this.props.onSave();
  };
  render() {
    let controls;
    if (this.props.isEditable) {
      controls = (
        <div>
          <p>
            <button className="usa-button-secondary editable-panel-cancel" onClick={this.handleCancelClick}>
              Cancel
            </button>
            <button
              className="usa-button editable-panel-save"
              onClick={this.handleSaveClick}
              disabled={!this.props.isValid}
            >
              Save
            </button>
          </p>
        </div>
      );
    }
    const classes = classNames(
      'editable-panel',
      {
        'is-editable': this.props.isEditable,
      },
      this.props.className,
    );
    return (
      <div className={classes}>
        <div className="editable-panel-header">
          <div className="title">{this.props.title}</div>
          {!this.props.isEditable &&
            this.props.editEnabled && (
              <a className="editable-panel-edit" onClick={this.handleEditClick}>
                Edit
              </a>
            )}
        </div>
        <div className="editable-panel-content">
          {this.props.children}
          {controls}
        </div>
      </div>
    );
  }
}
// TODO: Refactor when we switch to using a wizard
// Editable panel specific to Enter pre-move survey. Due to not using a wizard to enter the pre-move survey this
// panel has highly specific behavior (opening the edit view via clicking on Enter pre-move survey button)
export function PreMoveSurveyEditablePanelify(DisplayComponent, EditComponent, editEnabled = true) {
  const Wrapper = class extends Component {
    state = {
      isEditable: false,
    };
    componentDidUpdate = (prevProps, prevState) => {
      if (!prevProps.editPreMoveSurvey && this.props.editPreMoveSurvey) {
        this.setIsEditable(true);
      }
    };
    save = () => {
      let isValid = this.props.valid;
      if (isValid) {
        let args = this.props.getUpdateArgs();
        this.props.update(...args);
        this.setIsEditable(false);
      }
    };
    cancel = () => {
      this.props.reset();
      this.setIsEditable(false);
    };
    setIsEditable = isEditable => this.setState({ isEditable });
    render() {
      const isEditable = (editEnabled && (this.state.isEditable || this.props.isUpdating)) || false;
      const Content = isEditable ? EditComponent : DisplayComponent;
      return (
        <React.Fragment>
          {this.props.hasError && (
            <Alert type="error" heading="An error occurred">
              <em>{this.props.errorMessage}</em>
            </Alert>
          )}
          <PreMoveSurveyEditablePanel
            title={this.props.title}
            className={this.props.className}
            onSave={this.save}
            onEdit={this.setIsEditable}
            onCancel={this.cancel}
            isEditable={isEditable}
            editEnabled={editEnabled}
            isValid={this.props.valid}
          >
            <Content {...this.props} />
          </PreMoveSurveyEditablePanel>
        </React.Fragment>
      );
    }
  };
  Wrapper.propTypes = {
    update: PropTypes.func.isRequired,
    title: PropTypes.string.isRequired,
    isUpdating: PropTypes.bool,
  };
  return Wrapper;
}

PreMoveSurveyEditablePanel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  isEditable: PropTypes.bool.isRequired,
  editEnabled: PropTypes.bool,
  isValid: PropTypes.bool,
  onCancel: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};

PreMoveSurveyEditablePanel.defaultProps = {
  editEnabled: true,
};

const SurveyDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-3-column">
        <PanelSwaggerField title="Planned Pack Date" fieldName="pm_survey_planned_pack_date" required {...fieldProps} />
        <PanelSwaggerField
          title="Planned Pickup Date"
          fieldName="pm_survey_planned_pickup_date"
          required
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Planned Delivery Date"
          fieldName="pm_survey_planned_delivery_date"
          required
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-3-column">
        <PanelSwaggerField title="Weight Estimate" fieldName="pm_survey_weight_estimate" required {...fieldProps} />
        <PanelSwaggerField
          title="Progear Weight Estimate"
          fieldName="pm_survey_progear_weight_estimate"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Spouse Progear Weight Estimate"
          fieldName="pm_survey_spouse_progear_weight_estimate"
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-3-column">
        <PanelSwaggerField
          title="Pre-move survey conducted"
          fieldName="pm_survey_conducted_date"
          required
          {...fieldProps}
        />
        <PanelSwaggerField title="Survey Method" fieldName="pm_survey_method" required {...fieldProps} />
        <PanelSwaggerField title="Notes" fieldName="pm_survey_notes" className="notes" {...fieldProps} />
      </div>
    </React.Fragment>
  );
};

const SurveyEdit = props => {
  const schema = props.shipmentSchema;
  return (
    <React.Fragment>
      <FormSection name="survey">
        <div className="editable-panel-column">
          <SwaggerField fieldName="pm_survey_planned_pack_date" swagger={schema} required />
          <SwaggerField fieldName="pm_survey_planned_pickup_date" swagger={schema} required />
          <SwaggerField fieldName="pm_survey_planned_delivery_date" swagger={schema} required />
        </div>

        <div className="editable-panel-column">
          <SwaggerField fieldName="pm_survey_weight_estimate" swagger={schema} required />
          <SwaggerField fieldName="pm_survey_progear_weight_estimate" swagger={schema} />
          <SwaggerField fieldName="pm_survey_spouse_progear_weight_estimate" swagger={schema} />
        </div>
        <SwaggerField
          fieldName="pm_survey_conducted_date"
          title="Pre-move survey conducted"
          swagger={schema}
          required
        />
        <SwaggerField fieldName="pm_survey_method" swagger={schema} required />
        <SwaggerField fieldName="pm_survey_notes" title="Notes about dates" swagger={schema} />
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'shipment_pre_move_survey';

let PremoveSurveyPanel = PreMoveSurveyEditablePanelify(SurveyDisplay, SurveyEdit);

PremoveSurveyPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PremoveSurveyPanel);

PremoveSurveyPanel.propTypes = {
  shipment: PropTypes.object,
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  const updateShipmentStatus = getRequestStatus(state, premoveSurveyUpdateShipmentLabel);
  let hasError = false;
  let errorMessage = '';

  if (updateShipmentStatus.isLoading === false && updateShipmentStatus.isSuccess === false) {
    const errors = get(updateShipmentStatus, 'error.response.response.body.errors', {});
    errorMessage = humanReadableError(errors);
    hasError = true;
  }

  return {
    // reduxForm
    formValues: formValues,
    initialValues: {
      survey: pick(props.shipment, surveyFields),
    },

    shipmentSchema: get(state, 'swaggerPublic.spec.definitions.Shipment', {}),

    hasError,
    errorMessage,
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues.survey, premoveSurveyUpdateShipmentLabel];
    },
  };
}

export default connect(mapStateToProps)(PremoveSurveyPanel);
