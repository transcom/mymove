import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class PlaceInSitForm extends Component {
  render() {
    const { storageInTransitSchema, minDate, onSubmit } = this.props;
    return (
      <form onSubmit={this.props.handleSubmit(onSubmit)} className="place-in-sit-form">
        <div className="editable-panel-column">
          <SwaggerField
            className="place-in-sit-field"
            fieldName="actual_start_date"
            swagger={storageInTransitSchema}
            title="Actual start date"
            onChange={this.onChange}
            minDate={minDate}
            required
          />
        </div>
      </form>
    );
  }
}

PlaceInSitForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export const formName = 'place_in_sit_form';
PlaceInSitForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PlaceInSitForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransit', {}),
  };
}

export default connect(mapStateToProps)(PlaceInSitForm);
