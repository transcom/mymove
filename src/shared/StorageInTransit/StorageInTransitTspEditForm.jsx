import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

import './StorageInTransit.css';

export class StorageInTransitTspEditForm extends Component {
  render() {
    const { storageInTransitSchema, minDate } = this.props;
    const schema = this.props.storageInTransitSchema;
    const storageInTransit = this.props.initialValues;
    const fieldProps = {
      schema,
      values: this.props.initialValues,
    };

    const minActualStartDate = new Date(minDate);
    const utcMinDate = new Date(
      minActualStartDate.getUTCFullYear(),
      minActualStartDate.getUTCMonth(),
      minActualStartDate.getUTCDate(),
    );
    const maxActualStartDate = new Date(storageInTransit.out_date);
    const utcMaxDate = new Date(
      maxActualStartDate.getUTCFullYear(),
      maxActualStartDate.getUTCMonth(),
      maxActualStartDate.getUTCDate(),
    );
    const disabledDaysForDayPicker = [{ before: utcMinDate }, { after: utcMaxDate }];
    const inSit = storageInTransit.status === 'IN_SIT';
    const isReleased = storageInTransit.status === 'RELEASED';
    const isDelivered = storageInTransit.status === 'DELIVERED';
    const outDateMinDate = new Date(storageInTransit.actual_start_date);
    const utcOutDate = new Date(
      outDateMinDate.getUTCFullYear(),
      outDateMinDate.getUTCMonth(),
      outDateMinDate.getUTCDate(),
    );
    const disabledDaysForOutDayPicker = [{ before: utcOutDate }];

    return (
      <form onSubmit={this.props.handleSubmit(this.props.onSubmit)} className="storage-in-transit-tsp-edit-form">
        <div className="editable-panel-column">
          <PanelSwaggerField fieldName="location" required title="SIT location" {...fieldProps} />
          <PanelSwaggerField fieldName="estimated_start_date" required title="Estimated start date" {...fieldProps} />
          {inSit || isReleased || isDelivered ? (
            <SwaggerField
              className="storage-in-transit-form"
              fieldName="actual_start_date"
              swagger={storageInTransitSchema}
              title="Actual start date"
              minDate={minDate}
              disabledDays={disabledDaysForDayPicker}
              required
            />
          ) : (
            <PanelField
              value={storageInTransit.actual_start_date}
              fieldName="actual_start_date"
              required
              title="Actual start date"
              {...fieldProps}
            />
          )}
          {isReleased || isDelivered ? (
            <SwaggerField
              className="storage-in-transit-form"
              fieldName="out_date"
              disabledDays={disabledDaysForOutDayPicker}
              minDate={minDate}
              optional
              title="Date out"
              swagger={storageInTransitSchema}
            />
          ) : (
            <PanelField
              value={storageInTransit.out_date || 'n/a'}
              fieldName="days_out"
              required
              title="Days out"
              {...fieldProps}
            />
          )}
          <PanelField value="n/a" fieldName="days_used" required title="Days used" {...fieldProps} />
          <PanelField value="n/a" fieldName="expires" required title="Expires" {...fieldProps} />
          <SwaggerField fieldName="notes" optional title="Note" swagger={storageInTransitSchema} />
        </div>
        <div className="editable-panel-column">
          <div className="panel-subhead">Authorization</div>
          <PanelField
            value={storageInTransit.authorized_start_date}
            fieldName="authorized_start_date"
            title="Earliest authorized start"
            required
            {...fieldProps}
          />
          <PanelField
            value={storageInTransit.authorization_notes || 'n/a'}
            className="sit-approval-field"
            fieldName="authorization_notes"
            title="Note from reviewer"
            {...fieldProps}
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

StorageInTransitTspEditForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'storage_in_transit_tsp_edit_form';

StorageInTransitTspEditForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(StorageInTransitTspEditForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitTspEditForm);
