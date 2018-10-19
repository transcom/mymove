import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, pick } from 'lodash';
import { reduxForm, FormSection, getFormValues } from 'redux-form';

import { PanelSwaggerField, PanelField, editablePanelify } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { formatDate } from 'shared/formatters';

const datesFields = [
  'pm_survey_conducted_date',
  'pm_survey_planned_pack_date',
  'pm_survey_planned_pickup_date',
  'pm_survey_planned_delivery_date',
  'requested_pickup_date',
  'actual_pack_date',
  'actual_pickup_date',
  'actual_delivery_date',
  'requested_delivery_date',
  'pm_survey_notes',
  'pm_survey_method',
];

const DatesDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  // RDD is the best known date, so it prefers actual over planned over original.
  // TODO: cgilmer 2018/10/15 add `|| props.shipment.original_delivery_date` to this line
  const rdd = props.shipment.actual_delivery_date || props.shipment.pm_survey_planned_delivery_date;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <div className="column-subhead">PM Survey</div>
        <PanelSwaggerField title="PM survey conducted" fieldName="pm_survey_conducted_date" required {...fieldProps} />
        <PanelSwaggerField title="Survey Method" fieldName="pm_survey_method" required {...fieldProps} />
        <div className="column-subhead">Packing</div>
        <PanelField className="Todo-phase2 original_pack_date" title="Original" value="TODO" />
        <PanelSwaggerField fieldName="pm_survey_planned_pack_date" required title="Planned" {...fieldProps} />
        <PanelSwaggerField fieldName="actual_pack_date" required title="Actual" {...fieldProps} />
      </div>
      <div className="editable-panel-column">
        <div className="column-subhead">Pickup</div>
        <PanelSwaggerField fieldName="requested_pickup_date" required title="Original" {...fieldProps} />
        <PanelSwaggerField fieldName="pm_survey_planned_pickup_date" required title="Planned" {...fieldProps} />
        <PanelSwaggerField fieldName="actual_pickup_date" required title="Actual" {...fieldProps} />
        <div className="column-subhead">Delivery</div>
        <PanelField className="Todo-phase2 original_delivery_date" title="Original" value="TODO" />
        <PanelSwaggerField fieldName="pm_survey_planned_delivery_date" required title="Planned" {...fieldProps} />
        <PanelSwaggerField fieldName="actual_delivery_date" required title="Actual" {...fieldProps} />
        {/* TODO: cgilmer 2018/10/15 - replace ORIGINAL with original_delivery_date */}
        <PanelField className="rdd" title="Current RDD" value={(rdd && formatDate(rdd)) || 'ORIGINAL'} />
        <PanelSwaggerField fieldName="pm_survey_notes" title="Notes" {...fieldProps} />
      </div>
    </Fragment>
  );
};

const DatesEdit = props => {
  const schema = props.shipmentSchema;
  const fieldProps = {
    schema,
    values: props.shipment,
  };
  return (
    <Fragment>
      <FormSection name="dates">
        <div className="editable-panel-column">
          <div className="column-head">PM Survey</div>
          <SwaggerField fieldName="pm_survey_conducted_date" swagger={schema} />
          <SwaggerField fieldName="pm_survey_method" swagger={schema} />
          <div className="column-head">Packing</div>
          <PanelField title="Original" value="TODO" />
          <SwaggerField fieldName="pm_survey_planned_pack_date" title="Planned" swagger={schema} />
          <SwaggerField fieldName="actual_pack_date" title="Actual" swagger={schema} />
        </div>
        <div className="editable-panel-column">
          <div className="column-head">Pickup</div>
          <PanelSwaggerField fieldName="requested_pickup_date" title="Original" {...fieldProps} />
          <SwaggerField fieldName="pm_survey_planned_pickup_date" title="Planned" swagger={schema} />
          <SwaggerField fieldName="actual_pickup_date" title="Actual" swagger={schema} />
          <div className="column-head">Delivery</div>
          <PanelField title="Original" value="TODO" />
          <SwaggerField fieldName="pm_survey_planned_delivery_date" title="Planned" swagger={schema} />
          <SwaggerField fieldName="actual_delivery_date" title="Actual" swagger={schema} />
          <PanelField title="Current RDD" value="TODO" />
          <SwaggerField fieldName="pm_survey_notes" title="Notes about dates" swagger={schema} />
        </div>
      </FormSection>
    </Fragment>
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

    shipmentSchema: get(state, 'swaggerPublic.spec.definitions.Shipment', {}),

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
