import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { reduxForm, FormSection, getFormValues } from 'redux-form';

import {
  PanelSwaggerField,
  PanelField,
  editablePanelify,
} from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const datesFields = [
  'pm_survey_conducted_date',
  'pm_survey_planned_pack_date',
  'pm_survey_planned_pickup_date',
  'pm_survey_planned_delivery_date',
  'requested_pickup_date',
  'actual_pickup_date',
  'actual_pack_date',
  'requested_delivery_date',
  'actual_delivery_date',
  'pm_survey_notes',
  'pm_survey_method',
];

const DatesDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <div className="column-subhead">PM Survey</div>
        <PanelSwaggerField
          title="PM survey conducted"
          fieldName="pm_survey_conducted_date"
          required
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Survey Method"
          fieldName="pm_survey_method"
          required
          {...fieldProps}
        />
        <div className="column-subhead">Packing</div>
        <PanelField title="Original" value="TODO" />
        <PanelSwaggerField
          fieldName="pm_survey_planned_pack_date"
          required
          title="Planned"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="actual_pack_date"
          required
          title="Actual"
          {...fieldProps}
        />
      </div>
      <div className="editable-panel-column">
        <div className="column-subhead">Pickup</div>
        <PanelSwaggerField
          fieldName="requested_pickup_date"
          required
          title="Original"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="pm_survey_planned_pickup_date"
          required
          title="Planned"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="actual_pickup_date"
          required
          title="Actual"
          {...fieldProps}
        />
        <div className="column-subhead">Delivery</div>
        <PanelField title="Original" value="TODO" />
        <PanelSwaggerField
          fieldName="pm_survey_planned_delivery_date"
          required
          title="Planned"
          {...fieldProps}
        />
        <PanelField title="Current RDD" value="TODO" />
        <PanelSwaggerField
          fieldName="pm_survey_notes"
          required
          title="Notes about dates"
          {...fieldProps}
        />
      </div>
    </React.Fragment>
  );
};

const DatesEdit = props => {
  const schema = props.shipmentSchema;
  const fieldProps = {
    schema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <FormSection name="dates">
        <div className="editable-panel-column">
          <div className="column-head">PM Survey</div>
          <SwaggerField
            fieldName="pm_survey_conducted_date"
            swagger={schema}
            required
          />
          <SwaggerField
            fieldName="pm_survey_method"
            swagger={schema}
            required
          />
          <div className="column-head">Packing</div>
          <PanelField title="Original" value="TODO" />
          <SwaggerField
            fieldName="pm_survey_planned_pack_date"
            required
            title="Planned"
            swagger={schema}
          />
          <SwaggerField
            fieldName="actual_pack_date"
            required
            title="Actual"
            swagger={schema}
          />
        </div>
        <div className="editable-panel-column">
          <div className="column-head">Pickup</div>
          <PanelSwaggerField
            fieldName="requested_pickup_date"
            required
            title="Original"
            {...fieldProps}
          />
          <SwaggerField
            fieldName="pm_survey_planned_pickup_date"
            required
            title="Planned"
            swagger={schema}
          />
          <SwaggerField
            fieldName="actual_pickup_date"
            required
            title="Actual"
            swagger={schema}
          />
          <div className="column-head">Delivery</div>
          <PanelField title="Original" value="TODO" />
          <SwaggerField
            fieldName="pm_survey_planned_delivery_date"
            required
            title="Planned"
            swagger={schema}
          />
          <PanelField title="Current RDD" value="TODO" />
          <SwaggerField
            fieldName="pm_survey_notes"
            title="Notes about dates"
            swagger={schema}
          />
        </div>
      </FormSection>
    </React.Fragment>
  );
};

const formName = 'shipment_dates';

let DatesPanel = editablePanelify(DatesDisplay, DatesEdit);
DatesPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(DatesPanel);

DatesPanel.propTypes = {
  shipment: PropTypes.object,
  schema: PropTypes.object,
};

function mapStateToProps(state, props) {
  const formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues,
    initialValues: {
      dates: pick(props.shipment, datesFields),
    },

    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),

    hasError: !!props.error,
    errorMessage: props.error,
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues.dates];
    },
  };
}

export default connect(mapStateToProps)(DatesPanel);
