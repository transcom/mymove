import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get } from 'lodash';

import { push } from 'react-router-redux';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { updateServiceMember } from 'scenes/ServiceMembers/ducks';

import 'scenes/ServiceMembers/ServiceMembers.css';
import './Review.css';

const editContactFormName = 'edit_contact_info';

let EditContactForm = props => {
  const { onCancel, schema, handleSubmit, submitting, valid } = props;
  const resAddressSchema = get(schema, 'properties.residential_address');
  const backupAddressSchema = get(schema, 'properties.backup_mailing_address');
  return (
    <form className="service_member_contact_info" onSubmit={handleSubmit}>
      <h2>Edit Contact Info</h2>
      <SwaggerField fieldName="telephone" swagger={schema} required />
      <SwaggerField fieldName="secondary_telephone" swagger={schema} />
      <SwaggerField fieldName="personal_email" swagger={schema} required />
      <fieldset key="contact_preferences">
        <legend htmlFor="contact_preferences">
          Preferred contact method(s) during your move:
        </legend>
        <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
        <SwaggerField fieldName="text_message_is_preferred" swagger={schema} />
        <SwaggerField fieldName="email_is_preferred" swagger={schema} />
      </fieldset>
      <hr className="spacer" />
      <h3>Current Residence Address</h3>
      <SwaggerField
        fieldName="street_address_1"
        swagger={resAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="street_address_2"
        swagger={resAddressSchema || schema}
      />
      <SwaggerField
        fieldName="city"
        swagger={resAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="state"
        swagger={resAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="postal_code"
        swagger={resAddressSchema || schema}
        required
      />
      <hr className="spacer" />
      <h3>Backup Mailing Address</h3>
      <SwaggerField
        fieldName="street_address_1"
        swagger={backupAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="street_address_2"
        swagger={backupAddressSchema || schema}
      />
      <SwaggerField
        fieldName="city"
        swagger={backupAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="state"
        swagger={backupAddressSchema || schema}
        required
      />
      <SwaggerField
        fieldName="postal_code"
        swagger={backupAddressSchema || schema}
        required
      />
      <button type="submit" disabled={submitting || !valid}>
        Save
      </button>
      <button type="button" disabled={submitting} onClick={onCancel}>
        Cancel
      </button>
    </form>
  );
};

EditContactForm = reduxForm({
  form: editContactFormName,
})(EditContactForm);

class EditContact extends Component {
  returnToReview = () => {
    const reviewAddress = `/moves/${this.props.match.params.moveId}/review`;
    this.props.push(reviewAddress);
  };

  updateContact = fieldValues => {
    return this.props.updateServiceMember(fieldValues).then(() => {
      // This promise resolves regardless of error.
      if (!this.props.hasSubmitError) {
        this.returnToReview();
      } else {
        window.scrollTo(0, 0);
      }
    });
  };
  render() {
    const { schema, serviceMember } = this.props;
    var fullSM;
    // redux form expects a flat object for initial values to populate correctly
    // TODO: not working for some reason.
    if (
      get(serviceMember, 'residential_address') &&
      get(serviceMember, 'backup_mailing_address')
    ) {
      fullSM = Object.assign(
        {},
        serviceMember,
        serviceMember.residential_address,
        serviceMember.backup_mailing_address,
      );
    }
    return (
      <div className="usa-grid">
        <div className="usa-width-one-whole">
          <EditContactForm
            initialValues={fullSM}
            schema={schema}
            onSubmit={this.updateContact}
            onCancel={this.returnToReview}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    serviceMember: get(state, 'loggedInUser.loggedInUser.service_member'),
    move: get(state, 'moves.currentMove'),
    schema: get(
      state,
      'swagger.spec.definitions.CreateServiceMemberPayload',
      {},
    ),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ push, updateServiceMember }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(EditContact);
