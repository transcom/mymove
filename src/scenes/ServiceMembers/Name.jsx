import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

import SectionWrapper from 'components/Customer/SectionWrapper';

const subsetOfFields = ['first_name', 'middle_name', 'last_name', 'suffix'];
const formName = 'service_member_name';
const NameWizardForm = reduxifyWizardForm(formName);

export class Name extends Component {
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
        first_name: values.first_name,
        middle_name: values.middle_name,
        last_name: values.last_name,
        suffix: values.suffix,
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
    const { pages, pageKey, error, currentServiceMember } = this.props;
    const { errorMessage } = this.state;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <NameWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error || errorMessage}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <h1>Name</h1>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <SwaggerField fieldName="first_name" swagger={this.props.schema} required />
          </div>
          <SwaggerField fieldName="middle_name" swagger={this.props.schema} />
          <SwaggerField fieldName="last_name" swagger={this.props.schema} required />
          <SwaggerField fieldName="suffix" swagger={this.props.schema} />
        </SectionWrapper>
      </NameWizardForm>
    );
  }
}
Name.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    // TODO
    ...state.serviceMember,
    //
    currentServiceMember: serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(Name);
