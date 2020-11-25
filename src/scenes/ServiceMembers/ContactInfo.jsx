import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { selectCurrentUser } from 'shared/Data/users';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

const subsetOfFields = [
  'telephone',
  'secondary_telephone',
  'personal_email',
  'phone_is_preferred',
  'email_is_preferred',
];

const validateContactForm = (values) => {
  let errors = {};

  let prefSelected = Boolean(values.phone_is_preferred || values.email_is_preferred);
  if (!prefSelected) {
    const newError = {
      phone_is_preferred: 'Please select a preferred method of contact.',
    };
    return newError;
  }
  return errors;
};

const formName = 'service_member_contact_info';
const ContactWizardForm = reduxifyWizardForm(formName, validateContactForm);

export class ContactInfo extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  handleSubmit = () => {
    const { values, currentServiceMember, updateServiceMember } = this.props;
    if (values) {
      const payload = {
        id: currentServiceMember.id,
        ...values,
      };

      return patchServiceMember(payload)
        .then((response) => {
          updateServiceMember(response);
        })
        .catch((e) => {
          // TODO - error handling - below is rudimentary error handling to approximate existing UX
          // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
          const { response } = e;
          const errorMessage = getResponseError(response, 'failed to update service member due to server error');
          this.setState({
            errorMessage,
          });
        });
    }

    return Promise.resolve();
  };

  render() {
    const { pages, pageKey, error, currentServiceMember, userEmail, schema } = this.props;
    const { errorMessage } = this.state;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;
    if (initialValues && !initialValues.personal_email) {
      initialValues.personal_email = userEmail;
    }
    return (
      <ContactWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error || errorMessage}
        initialValues={initialValues}
      >
        <h1>Your contact info</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <SwaggerField fieldName="telephone" swagger={schema} required />
          </div>
          <SwaggerField fieldName="secondary_telephone" swagger={schema} />
          <SwaggerField fieldName="personal_email" swagger={schema} required />
          <fieldset className="usa-fieldset" key="contact_preferences">
            <p htmlFor="contact_preferences">Preferred contact method(s) during your move:</p>
            <SwaggerField fieldName="phone_is_preferred" swagger={schema} />
            <SwaggerField fieldName="email_is_preferred" swagger={schema} />
          </fieldset>
        </SectionWrapper>
      </ContactWizardForm>
    );
  }
}
ContactInfo.propTypes = {
  userEmail: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const user = selectCurrentUser(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    userEmail: user.email,
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    // TODO
    ...state.serviceMember,
    //
    currentServiceMember: serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ContactInfo);
