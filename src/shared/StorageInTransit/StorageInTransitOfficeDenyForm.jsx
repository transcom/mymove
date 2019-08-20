import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class StorageInTransitOfficeDenyForm extends Component {
  render() {
    const { storageInTransitSchema } = this.props;
    return (
      <form onSubmit={this.props.handleSubmit(this.props.onSubmit)} className="storage-in-transit-office-deny-form">
        <fieldset key="sit-deny-information">
          <div className="editable-panel-column">
            <SwaggerField
              fieldName="authorization_notes"
              swagger={storageInTransitSchema}
              title="Reason for denial"
              required
            />
          </div>
        </fieldset>
      </form>
    );
  }
}

StorageInTransitOfficeDenyForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
};

export const formName = 'storage_in_transit_office_deny_form';
StorageInTransitOfficeDenyForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(StorageInTransitOfficeDenyForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
  };
}

export default connect(mapStateToProps)(StorageInTransitOfficeDenyForm);
