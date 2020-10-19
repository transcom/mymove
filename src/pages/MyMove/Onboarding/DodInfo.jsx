/* eslint-disable camelcase */
/* eslint-disable react/forbid-prop-types */
import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { getFormValues, Field } from 'redux-form';

import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { normalizeSSN } from 'shared/JsonSchemaForm/reduxFieldNormalizer';
import SSNField from 'components/form/fields/SSNInput';
import { patchServiceMember } from 'services/internalApi';

const validateDodForm = (values, form) => {
  // Everything is taken care of except for SSN
  const errors = {};
  const ssn = values.social_security_number;
  const hasSSN = form.ssnOnServer;

  const validSSNPattern = RegExp('^\\d{3}-\\d{2}-\\d{4}$');
  const validSSN = validSSNPattern.test(ssn);
  const ssnPresent = ssn !== '' && ssn !== undefined;

  if (hasSSN) {
    if (ssnPresent && !validSSN) {
      errors.social_security_number = 'SSN must have 9 digits';
    }
  } else if (!ssnPresent) {
    errors.social_security_number = 'Required';
  } else if (!validSSN) {
    errors.social_security_number = 'SSN must have 9 digits';
  }

  return errors;
};

const formName = 'service_member_dod_info';
const DodWizardForm = reduxifyWizardForm(formName, validateDodForm);

export class DodInfo extends Component {
  handleSubmit = () => {
    const { values, currentServiceMember, updateServiceMember } = this.props;
    if (values) {
      const patchServiceMemberPayload = {
        id: currentServiceMember.id,
        affiliation: values.affiliation,
        edipi: values.edipi,
        social_security_number: values.social_security_number,
        rank: values.rank,
      };

      return patchServiceMember(patchServiceMemberPayload)
        .then((response) => {
          updateServiceMember(response);
        })
        .catch(() => {
          // TODO
          // console.log('handle errors inline', e);
        });
    }

    return Promise.resolve();
  };

  render() {
    const { pages, pageKey, error, currentServiceMember, schema } = this.props;
    const { affiliation, edipi, social_security_number, rank } = currentServiceMember;

    const initialValues = {
      affiliation,
      edipi,
      social_security_number,
      rank,
    };

    const ssnOnServer = currentServiceMember?.has_social_security_number || false;

    return (
      <DodWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
        ssnOnServer={ssnOnServer}
      >
        <div className="grid-row">
          <div className="grid-col-12">
            <h1 className="sm-heading">Create your profile</h1>
            <p>Before we can schedule your move, we need to know a little more about you.</p>
            <SwaggerField fieldName="affiliation" swagger={schema} required />
            <SwaggerField fieldName="edipi" swagger={schema} required />
            <Field
              name="social_security_number"
              component={SSNField}
              ssnOnServer={ssnOnServer}
              normalize={normalizeSSN}
            />
            <SwaggerField fieldName="rank" swagger={schema} required />
          </div>
        </div>
      </DodWizardForm>
    );
  }
}

DodInfo.propTypes = {
  pageKey: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  values: PropTypes.object,
  pages: PropTypes.arrayOf(PropTypes.string).isRequired,
  updateServiceMember: PropTypes.func,
};

DodInfo.defaultProps = {
  currentServiceMember: null,
  error: null,
  values: null,
  updateServiceMember: () => {},
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

const mapStateToProps = (state) => {
  const props = {
    schema: state.swaggerInternal?.spec?.definitions?.CreateServiceMemberPayload || {},
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
  return props;
};

export default connect(mapStateToProps, mapDispatchToProps)(DodInfo);
