import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';

const subsetOfFields = ['affiliation', 'edipi', 'rank'];

const formName = 'service_member_dod_info';
const DodWizardForm = reduxifyWizardForm(formName);

export class DodInfo extends Component {
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
        affiliation: values.affiliation,
        edipi: values.edipi,
        rank: values.rank,
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
    const { pages, pageKey, error, currentServiceMember, schema } = this.props;
    const { errorMessage } = this.state;

    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;

    return (
      <DodWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error || errorMessage}
        initialValues={initialValues}
      >
        <h1>Create your profile</h1>
        <p>Before we can schedule your move, we need to know a little more about you.</p>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <SwaggerField fieldName="affiliation" swagger={schema} required />
          </div>
          <SwaggerField fieldName="edipi" swagger={schema} required />
          <SwaggerField fieldName="rank" swagger={schema} required />
        </SectionWrapper>
      </DodWizardForm>
    );
  }
}

DodInfo.propTypes = {
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

  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    // TODO
    ...state.serviceMember,
    //
    currentServiceMember: serviceMember,
  };
  return props;
}

export default connect(mapStateToProps, mapDispatchToProps)(DodInfo);
