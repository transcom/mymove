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

const weightsFields = ['actual_weight'];

const DatesDisplay = props => {
  const fieldProps = {
    schema: props.shipmentSchema,
    values: props.shipment,
  };
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <div className="column-subhead">PM Survey</div>
        <PanelSwaggerField fieldName="actual_weight" required {...fieldProps} />
        <div className="column-subhead">Packing</div>
        <PanelSwaggerField
          fieldName="pm_survey_planned_pack_date"
          required
          title="Planned"
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
        <PanelSwaggerField
          fieldName="pm_survey_planned_delivery_date"
          required
          title="Planned"
          {...fieldProps}
        />
        <PanelSwaggerField
          fieldName="actual_delivery_date"
          required
          title="Actual"
          {...fieldProps}
        />
        <PanelField title="RDD" value="RDD TK" />
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
          <div className="column-subhead">PM Survey</div>

          <SwaggerField fieldName="actual_weight" swagger={schema} />
          <div className="column-subhead">Packing</div>
          <SwaggerField
            fieldName="pm_survey_planned_pack_date"
            required
            title="Planned"
            swagger={schema}
          />
        </div>
        <div className="editable-panel-column">
          <div className="column-subhead">Pickup</div>
          <SwaggerField
            fieldName="requested_pickup_date"
            required
            title="Original"
            swagger={schema}
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
          <div className="column-subhead">Delivery</div>
          <SwaggerField
            fieldName="pm_survey_planned_delivery_date"
            required
            title="Planned"
            swagger={schema}
          />
          <SwaggerField
            fieldName="actual_delivery_date"
            required
            title="Actual"
            swagger={schema}
          />
          <PanelField title="RDD" value="RDD TK" />
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
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  return {
    // reduxForm
    formValues,
    initialValues: {
      shipment: props.shipment,
    },

    shipmentSchema: get(state, 'swagger.spec.definitions.Shipment', {}),

    hasError: !!props.error,
    errorMessage: props.error,
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues.shipment];
    },
  };
}

export default connect(mapStateToProps)(DatesPanel);
