import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { getFormValues } from 'redux-form';

import { patchServiceMember, getResponseError } from 'services/internalApi';
import { updateServiceMember as updateServiceMemberAction } from 'store/entities/actions';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import AddressForm from 'shared/AddressForm';

import SectionWrapper from 'components/Customer/SectionWrapper';

const formName = 'service_member_backup_mailing_addresss';
const BackupMailingWizardForm = reduxifyWizardForm(formName);

export class BackupMailingAddress extends Component {
  constructor(props) {
    super(props);

    this.state = {
      errorMessage: null,
    };
  }

  handleSubmit = () => {
    const { values, currentServiceMember, updateServiceMember } = this.props;

    const payload = {
      id: currentServiceMember.id,
      backup_mailing_address: values,
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
  };

  render() {
    const { pages, pageKey, error, currentServiceMember } = this.props;
    const { errorMessage } = this.state;

    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = get(currentServiceMember, 'backup_mailing_address');
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <BackupMailingWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error || errorMessage}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <h1>Backup mailing address</h1>
        <p>
          Where should we send mail if we can’t reach you at your primary address? You might use a parent's or friend’s
          address, or a post office box.
        </p>
        <SectionWrapper>
          <div className="tablet:margin-top-neg-3">
            <AddressForm schema={this.props.schema} />
          </div>
        </SectionWrapper>
      </BackupMailingWizardForm>
    );
  }
}
BackupMailingAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

const mapDispatchToProps = {
  updateServiceMember: updateServiceMemberAction,
};

function mapStateToProps(state) {
  return {
    schema: get(state, 'swaggerInternal.spec.definitions.Address', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(BackupMailingAddress);
