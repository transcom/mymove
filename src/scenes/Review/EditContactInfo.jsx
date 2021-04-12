import React, { Component } from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import scrollToTop from 'shared/scrollToTop';

import { push } from 'connected-react-router';
import { reduxForm, FormSection } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import Alert from 'shared/Alert';
import AddressForm from 'shared/AddressForm';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';

import './Review.css';
import SaveCancelButtons from './SaveCancelButtons';

const editContactFormName = 'edit_contact_info';

let EditContactForm = (props) => {
  const { serviceMemberSchema, addressSchema, handleSubmit, submitting, valid } = props;
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form className="service_member_contact_info" onSubmit={handleSubmit}>
            <FormSection name="serviceMember">
              <h1>Edit Contact Info</h1>
              <SwaggerField fieldName="telephone" swagger={serviceMemberSchema} required />
              <SwaggerField fieldName="secondary_telephone" swagger={serviceMemberSchema} />
              <SwaggerField fieldName="personal_email" swagger={serviceMemberSchema} required />
              <fieldset className="usa-fieldset" key="contact_preferences">
                <p htmlFor="contact_preferences">Preferred contact method(s) during your move:</p>
                <SwaggerField fieldName="phone_is_preferred" swagger={serviceMemberSchema} />
                <SwaggerField fieldName="email_is_preferred" swagger={serviceMemberSchema} />
              </fieldset>
            </FormSection>
            <hr className="spacer" />

            <FormSection name="resAddress">
              <h3>Current Residence Address</h3>
              <AddressForm schema={addressSchema} />
            </FormSection>
            <hr className="spacer" />
            <FormSection name="backupAddress">
              <h3>Backup Mailing Address</h3>
              <AddressForm schema={addressSchema} />
            </FormSection>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

const validateEditContactFormBools = (fields) => {
  return (values, form) => {
    let errors = {};
    let prefSelected = false;
    fields.forEach((fieldName) => {
      if (Boolean(get(values, fieldName))) {
        prefSelected = true;
      }
    });
    if (!prefSelected) {
      let valueSection = fields[0].split('.')[0];
      let field = fields[0].split('.')[1];
      var errorMsg = {
        [field]: 'Please select a preferred method of contact.',
      };
      var newError = { [valueSection]: errorMsg };
      return newError;
    }
    return errors;
  };
};

EditContactForm = reduxForm({
  form: editContactFormName,
  validate: validateEditContactFormBools(['serviceMember.phone_is_preferred', 'serviceMember.email_is_preferred']),
})(EditContactForm);

class EditContact extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  updateContact = (fieldValues) => {
    const { setFlashMessage } = this.props;

    let serviceMember = fieldValues.serviceMember;
    serviceMember.id = this.props.serviceMember.id;
    serviceMember.residential_address = fieldValues.resAddress;
    serviceMember.backup_mailing_address = fieldValues.backupAddress;

    return patchServiceMember(serviceMember)
      .then((response) => {
        // Update Redux with new data
        this.props.updateServiceMember(response);
        setFlashMessage('EDIT_CONTACT_INFO_SUCCESS', 'success', '', 'Your changes have been saved.');
        this.props.history.goBack();
      })
      .catch((e) => {
        // TODO - error handling - below is rudimentary error handling to approximate existing UX
        // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
        const { response } = e;
        const errorMessage = getResponseError(response, 'failed to update service member due to server error');
        this.setState({
          errorMessage,
        });

        scrollToTop();
      });
  };

  render() {
    const { serviceMemberSchema, addressSchema, serviceMember } = this.props;
    const { errorMessage } = this.state;

    let initialValues = null;
    if (serviceMember && get(serviceMember, 'residential_address') && get(serviceMember, 'backup_mailing_address'))
      initialValues = {
        serviceMember: serviceMember,
        resAddress: serviceMember.residential_address,
        backupAddress: serviceMember.backup_mailing_address,
      };

    return (
      <div className="usa-grid">
        {errorMessage && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {errorMessage}
            </Alert>
          </div>
        )}
        <div className="usa-width-one-whole">
          <EditContactForm
            initialValues={initialValues}
            serviceMemberSchema={serviceMemberSchema}
            addressSchema={addressSchema}
            onSubmit={this.updateContact}
          />
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    serviceMember,
    serviceMemberSchema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    addressSchema: get(state, 'swaggerInternal.spec.definitions.Address', {}),
  };
}

const mapDispatchToProps = {
  push,
  updateServiceMember: updateServiceMemberAction,
  setFlashMessage: setFlashMessageAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(EditContact);
