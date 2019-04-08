import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class PlaceInSitForm extends Component {
  //form submission still to be implemented
  handleSubmit = e => {
    e.preventDefault();
  };

  render() {
    const { storageInTransitSchema } = this.props;
    return (
      <form onSubmit={this.handleSubmit} className="place-in-sit-form">
        <div className="editable-panel-column">
          <SwaggerField
            className="place-in-sit-field"
            fieldName="actual_start_date"
            swagger={storageInTransitSchema}
            title="Actual start date"
            required
          />
        </div>
      </form>
    );
  }
}

PlaceInSitForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
};

const formName = 'place_in_sit_form';
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
