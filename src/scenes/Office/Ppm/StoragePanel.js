import { get } from 'lodash';
import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { reduxForm, getFormValues } from 'redux-form';

import { editablePanelify, PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectPPMForMove, updatePPM } from 'shared/Entities/modules/ppms';
import { formatCents } from '../../../shared/formatters';
import { convertDollarsToCents } from '../../../shared/utils';

const StorageDisplay = props => {
  const fieldProps = {
    schema: props.ppmSchema,
    values: {
      total_sit_cost: formatCents(props.ppm.total_sit_cost),
      days_in_storage: props.ppm.days_in_storage,
    },
  };

  return (
    <div className="editable-panel-column">
      <PanelSwaggerField fieldName="total_sit_cost" title="Total cost" {...fieldProps} />
      <PanelSwaggerField fieldName="days_in_storage" title="Total days in storage" {...fieldProps} />
    </div>
  );
};

const StorageEdit = props => {
  const schema = props.ppmSchema;

  return (
    <div className="editable-panel-column">
      <SwaggerField className="short-field storage" fieldName="total_sit_cost" title="Total cost" swagger={schema} />
      <SwaggerField
        className="short-field storage"
        fieldName="days_in_storage"
        title="Total days in storage"
        swagger={schema}
      />
    </div>
  );
};

const formName = 'ppm_sit_storage';

let StoragePanel = editablePanelify(StorageDisplay, StorageEdit);
StoragePanel = reduxForm({
  form: formName,
  enableReinitialize: true,
})(StoragePanel);

function mapStateToProps(state, props) {
  const formValues = getFormValues(formName)(state);
  const ppm = selectPPMForMove(state, props.moveId);

  return {
    // reduxForm
    formValues,
    initialValues: {
      total_sit_cost: ppm.total_sit_cost,
      days_in_storage: ppm.days_in_storage,
    },

    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,

    hasError: !!props.error,
    errorMessage: get(state, 'office.error'),
    isUpdating: false,

    // editablePanelify
    getUpdateArgs: function() {
      const values = getFormValues(formName)(state);
      const adjustedValues = {
        total_sit_cost: convertDollarsToCents(values.total_sit_cost),
        days_in_storage: values.days_in_storage,
      };
      return [props.moveId, ppm.id, adjustedValues];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updatePPM,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(StoragePanel);
