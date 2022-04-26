import { get } from 'lodash';
import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { reduxForm, getFormValues } from 'redux-form';
import { filter } from 'lodash';

import Alert from 'shared/Alert';
import { editablePanelify, PanelSwaggerField } from 'shared/EditablePanel';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectActivePPMForMove, updatePPM } from 'shared/Entities/modules/ppms';
import { formatCents } from 'utils/formatters';
import { convertDollarsToCents } from '../../../shared/utils';
import { getDocsByStatusAndType } from './ducks';

const StorageDisplay = (props) => {
  const cost = props.ppm && props.ppm.total_sit_cost ? formatCents(props.ppm.total_sit_cost) : 0;
  const days = props.ppm && props.ppm.days_in_storage ? props.ppm.days_in_storage : 0;

  const fieldProps = {
    schema: {
      properties: {
        days_in_storage: {
          maximum: 90,
          minimum: 0,
          title: 'How many days do you plan to put your stuff in storage?',
          type: 'integer',
          'x-nullable': true,
        },
        total_sit_cost: {
          type: 'string',
        },
      },
    },
    values: {
      total_sit_cost: `$${cost}`,
      days_in_storage: `${days}`,
    },
  };

  return (
    <div className="editable-panel-column">
      {props.awaitingStorageExpenses.length > 0 && (
        <div className="awaiting-storage-expenses-warning">
          <Alert type="warning">There are more storage receipts awaiting review</Alert>
        </div>
      )}
      <PanelSwaggerField fieldName="total_sit_cost" title="Total storage cost" {...fieldProps} />
      <PanelSwaggerField fieldName="days_in_storage" title="Days in storage" {...fieldProps} />
    </div>
  );
};

const StorageEdit = (props) => {
  const schema = props.ppmSchema;

  return (
    <div className="editable-panel-column">
      <SwaggerField
        className="short-field storage"
        fieldName="total_sit_cost"
        title="Total storage cost"
        swagger={schema}
      />
      <SwaggerField
        className="short-field storage"
        fieldName="days_in_storage"
        title="Days in storage"
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
  const ppm = selectActivePPMForMove(state, props.moveId);
  const storageExpenses = filter(props.moveDocuments, ['moving_expense_type', 'STORAGE']);

  return {
    // reduxForm
    formValues,
    initialValues: {
      total_sit_cost: formatCents(ppm.total_sit_cost),
      days_in_storage: ppm.days_in_storage,
    },

    ppmSchema: get(state, 'swaggerInternal.spec.definitions.PersonallyProcuredMovePayload'),
    ppm,

    hasError: !!props.error,
    errorMessage: get(state, 'office.error'),
    isUpdating: false,
    awaitingStorageExpenses: getDocsByStatusAndType(storageExpenses, 'OK'),

    // editablePanelify
    getUpdateArgs: function () {
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
