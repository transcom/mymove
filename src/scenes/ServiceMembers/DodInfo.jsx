import { get, pick } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { updateServiceMember } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import SectionWrapper from 'components/Customer/SectionWrapper';

const subsetOfFields = ['affiliation', 'edipi', 'rank'];

const formName = 'service_member_dod_info';
const DodWizardForm = reduxifyWizardForm(formName);

export class DodInfo extends Component {
  handleSubmit = () => {
    const pendingValues = this.props.values;
    if (pendingValues) {
      const patch = pick(pendingValues, subsetOfFields);
      return this.props.updateServiceMember(patch);
    }
  };

  render() {
    const { pages, pageKey, error, currentServiceMember, schema } = this.props;
    const initialValues = currentServiceMember ? pick(currentServiceMember, subsetOfFields) : null;

    return (
      <DodWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
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

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    schema: get(state, 'swaggerInternal.spec.definitions.CreateServiceMemberPayload', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(DodInfo);
