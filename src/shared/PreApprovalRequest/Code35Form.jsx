import React, { Component, Fragment } from 'react';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { Code35FormAlert } from './Code35FormAlert';

export class Code35Form extends Component {
  makeStaticFields() {
    return (
      <Fragment>
        <label htmlFor="description" className="usa-input-label">
          Description of service
        </label>
        <div>
          <strong>{this.props.initialValues.description}</strong>
        </div>
        <label htmlFor="reason" className="usa-input-label">
          Reason for service
        </label>
        <div>
          <strong>{this.props.initialValues.reason}</strong>
        </div>
        <label htmlFor="estimate_amount_cents" className="usa-input-label">
          Estimate, not to exceed
        </label>
        <div>
          <strong>{`$${this.props.initialValues.estimate_amount_cents}`}</strong>
        </div>
      </Fragment>
    );
  }

  makeEditableFields() {
    return (
      <Fragment>
        <SwaggerField
          className="textarea-half"
          title="Description of service"
          fieldName="description"
          swagger={this.props.ship_line_item_schema}
          required
        />
        <SwaggerField
          className="textarea-half"
          title="Reason for service"
          fieldName="reason"
          swagger={this.props.ship_line_item_schema}
          required
        />
        <SwaggerField
          title="Estimate, not to exceed"
          fieldName="estimate_amount_cents"
          swagger={this.props.ship_line_item_schema}
          required
        />
      </Fragment>
    );
  }

  render() {
    return (
      <Fragment>
        {this.props.status === 'CONDITIONALLY_APPROVED' || this.props.status === 'APPROVED'
          ? this.makeStaticFields()
          : this.makeEditableFields()}
        <SwaggerField
          title="Actual amount of service"
          fieldName="actual_amount_cents"
          swagger={this.props.ship_line_item_schema}
        />
        <div className="bq-explanation">
          <p>Enter amount after service is completed</p>
        </div>
        <Code35FormAlert showAlert={this.props.showAlert} />
      </Fragment>
    );
  }
}

Code35Form.propTypes = {};
