import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

import './StorageInTransit.css';

export class StorageInTransitOfficeEditForm extends Component {
  render() {
    const { storageInTransitSchema } = this.props;
    const schema = this.props.storageInTransitSchema;
    const storageInTransit = this.props.initialValues;
    const fieldProps = {
      schema,
      values: this.props.initialValues,
    };

    const sitStatus = () => {
      if (storageInTransit.status === 'APPROVED') {
        return 'Yes';
      } else {
        return 'No';
      }
    };

    return (
      <form onSubmit={this.props.handleSubmit(this.props.onSubmit)} className="storage-in-transit-form">
        <div className="editable-panel-column">
          <PanelSwaggerField fieldName="location" required title="SIT location" {...fieldProps} />
          <PanelSwaggerField fieldName="estimated_start_date" required title="Estimated start date" {...fieldProps} />
          <PanelField value="n/a" fieldName="actual_start_date" required title="Actual start date" {...fieldProps} />
          <PanelField value="n/a" fieldName="out_date" required title="Date out" {...fieldProps} />
          <PanelField value="n/a" fieldName="days_used" required title="Days used" {...fieldProps} />
          <PanelField value="n/a" fieldName="expires" required title="Expires" {...fieldProps} />
          <PanelSwaggerField fieldName="notes" optional title="Note" {...fieldProps} />
        </div>
        <div className="editable-panel-column">
          <div className="panel-subhead">Authorization</div>
          <PanelField value={sitStatus()} fieldName="status" required title="SIT Approved" {...fieldProps} />
          {storageInTransit.status === 'APPROVED' ? (
            <SwaggerField
              fieldName="authorized_start_date"
              swagger={storageInTransitSchema}
              title="Earliest authorized start"
              required
            />
          ) : (
            <PanelField
              value="n/a"
              fieldName="authorized_start_date"
              required
              title="Earliest authorized start"
              {...fieldProps}
            />
          )}
          <SwaggerField
            className="sit-approval-field"
            fieldName="authorization_notes"
            title="Note from reviewer"
            swagger={storageInTransitSchema}
          />
          <PanelField value="n/a" fieldName="sit_number" title="SIT number" {...fieldProps} />
          <div className="panel-field">
            <span className="field-value warehouse-field-margin">Warehouse</span>
            <span className="field-title">{storageInTransit.warehouse_name}</span>
            <span className="field-title">{storageInTransit.warehouse_address.street_address_1}</span>
            <span className="field-title">{storageInTransit.warehouse_address.street_address_2}</span>
            <span className="field-title">{storageInTransit.warehouse_address.street_address_3}</span>
            <span className="field-title">
              {storageInTransit.warehouse_address.city}, {storageInTransit.warehouse_address.state}{' '}
              {storageInTransit.warehouse_address.postal_code}
            </span>
          </div>
        </div>
      </form>
    );
  }
}

StorageInTransitOfficeEditForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'storage_in_transit_office_edit_form';

StorageInTransitOfficeEditForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(StorageInTransitOfficeEditForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitOfficeEditForm);
