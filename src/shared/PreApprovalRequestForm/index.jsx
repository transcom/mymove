import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { reduxForm, Field, Form } from 'redux-form';
import validator from 'shared/JsonSchemaForm/validator';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

class PreApprovalRequestForm extends Component {
  render() {
    return (
      <Form onSubmit={this.props.handleSubmit(this.props.onSubmit)}>
        <div className="rounded">
          <Field component="select" name="code" validate={validator.isRequired}>
            <option />
            {this.props.accessorials.map(e => (
              <option key={e.id} value={e.id}>
                {e.code} {e.item}
              </option>
            ))}
          </Field>
          <SwaggerField
            fieldName="location"
            swagger={this.props.ship_accessorial_schema}
            required
          />
          <SwaggerField
            fieldName="quantity_1"
            swagger={this.props.ship_accessorial_schema}
            required
          />
          <SwaggerField
            fieldName="notes"
            swagger={this.props.ship_accessorial_schema}
            required
          />
        </div>
      </Form>
    );
  }
}

PreApprovalRequestForm.propTypes = {
  schema: PropTypes.object,
  accessorials: PropTypes.array,
  onSubmit: PropTypes.func.isRequired,
};

const formName = 'preapproval_request_form';

PreApprovalRequestForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PreApprovalRequestForm);

function mapStateToProps(state, props) {
  return {
    // reduxForm
    initialValues: {
      location: 'ORIGIN',
    },
    accessorials: [
      {
        id: 'sdlfkj',
        code: 'F9D',
        item: 'Long Haul',
      },
    ],
    ship_accessorial_schema: get(
      state,
      'swagger.spec.definitions.ShipmentAccessorialPayload',
      {},
    ),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(
  PreApprovalRequestForm,
);
