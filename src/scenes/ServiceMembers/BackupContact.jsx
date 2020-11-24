import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { updateBackupContact as updateBackupContactAction } from 'store/entities/actions';
import { createBackupContactForServiceMember, patchBackupContact, getResponseError } from 'services/internalApi';
import { renderField, recursivelyAnnotateRequiredFields } from 'shared/JsonSchemaForm';
import { reduxForm } from 'redux-form';
import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import scrollToTop from 'shared/scrollToTop';

import SectionWrapper from 'components/Customer/SectionWrapper';

import './BackupContact.css';

const NonePermission = 'NONE';

const formName = 'service_member_backup_contact';

class ContactForm extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  static getDerivedStateFromProps(nextProps, prevState) {
    const { valid, dirty, updateValidDirty } = nextProps;
    updateValidDirty(valid, dirty);

    return prevState;
  }

  render() {
    const { schema } = this.props;
    recursivelyAnnotateRequiredFields(schema);
    const fields = schema.properties || {};

    return (
      <form>
        <h1>Backup contact</h1>
        <p>If we can't reach you, who can we contact (such as spouse or parent)?</p>
        <p>Any person you assign as a backup contact must be 18 years of age or older.</p>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            {renderField('name', fields, '')}
            {renderField('email', fields, '')}
            {renderField('telephone', fields, '')}
          </div>
        </SectionWrapper>
      </form>
    );
  }
}

const validateContact = (values, form) => {
  let requiredErrors = {};
  ['name', 'email'].forEach((requiredFieldName) => {
    if (values[`${requiredFieldName}`] === undefined || values[`${requiredFieldName}`] === '') {
      requiredErrors[`${requiredFieldName}`] = 'Required.';
    }
  });
  return requiredErrors;
};

ContactForm = reduxForm({ form: formName, validate: validateContact })(ContactForm);

export class BackupContact extends Component {
  constructor(props) {
    super(props);
    this.state = {
      isValid: true,
      isDirty: false,
      errorMessage: null,
    };
  }

  handleSubmit = () => {
    const { values, updateBackupContact, currentBackupContacts, match } = this.props;

    if (values) {
      const payload = {
        ...values,
        telephone: values.telephone === '' ? null : values.telephone,
        permission: values.permission === undefined ? NonePermission : values.permission,
      };

      if (currentBackupContacts.length > 0) {
        const [firstBackupContact] = currentBackupContacts;
        payload.id = firstBackupContact.id;
        return patchBackupContact(payload)
          .then((response) => {
            updateBackupContact(response);
          })
          .catch((e) => {
            // TODO - error handling - below is rudimentary error handling to approximate existing UX
            // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
            const { response } = e;
            const errorMessage = getResponseError(response, 'failed to update backup contact due to server error');
            this.setState({
              errorMessage,
            });

            scrollToTop();
          });
      } else {
        const { serviceMemberId } = match.params;
        return createBackupContactForServiceMember(serviceMemberId, payload)
          .then((response) => {
            updateBackupContact(response);
          })
          .catch((e) => {
            // TODO - error handling - below is rudimentary error handling to approximate existing UX
            // Error shape: https://github.com/swagger-api/swagger-js/blob/master/docs/usage/http-client.md#errors
            const { response } = e;
            const errorMessage = getResponseError(response, 'failed to create backup contact due to server error');
            this.setState({
              errorMessage,
            });

            scrollToTop();
          });
      }
    }

    return Promise.resolve();
  };

  updateValidDirty = (isValid, isDirty) => {
    this.setState({
      isValid,
      isDirty,
    });
  };

  render() {
    const { pages, pageKey, error } = this.props;
    const { isValid, isDirty, errorMessage } = this.state;

    //
    var [contact1, contact2] = this.props.currentBackupContacts; // contact2 will be used when we implement saving two backup contacts.

    // initialValues has to be null until there are values from the action since only the first values are taken
    const firstInitialValues = contact1 ? pick(contact1, ['name', 'email', 'telephone', 'permission']) : null;

    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        dirty={isDirty}
        error={error || errorMessage}
      >
        <ContactForm
          ref="currentForm"
          updateValidDirty={this.updateValidDirty}
          initialValues={firstInitialValues}
          handleSubmit={no_op}
          schema={this.props.schema}
        />
      </WizardPage>
    );
  }
}
BackupContact.propTypes = {
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

const mapDispatchToProps = {
  updateBackupContact: updateBackupContactAction,
};

function mapStateToProps(state) {
  return {
    currentBackupContacts: state.serviceMember.currentBackupContacts,
    error: state.serviceMember.error,
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberBackupContactPayload', {}),
    values: getFormValues(formName)(state),
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(BackupContact);
