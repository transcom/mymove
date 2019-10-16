import React, { Component } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import { ConnectedRouter } from 'react-router-redux';
import { history } from 'shared/store';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import Loadable from 'react-loadable';
import QueueHeader from 'shared/Header/Office';
import QueueList from './QueueList';
import QueueTable from './QueueTable';
import MoveInfo from './MoveInfo';
import OrdersInfo from './OrdersInfo';
import DocumentViewer from './DocumentViewer';
import { getCurrentUserInfo, selectCurrentUser } from 'shared/Data/users';
import { loadInternalSchema, loadPublicSchema } from 'shared/Swagger/ducks';
import { detectIE11, no_op } from 'shared/utils';
import LogoutOnInactivity from 'shared/User/LogoutOnInactivity';
import PrivateRoute from 'shared/User/PrivateRoute';
import ScratchPad from 'shared/ScratchPad';
import { isProduction } from 'shared/constants';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { RetrieveMovesForOffice } from './api';
import './office.scss';
import CustomerDetails from './TOO/customerDetails';
import { withContext } from 'shared/AppContext';

const TOO = Loadable({
  loader: () => import('./TOO/too'),
  loading: () => <LoadingPlaceholder />,
});

export class Queues extends Component {
  render() {
    return (
      <div className="usa-grid grid-wide queue-columns">
        <div className="queue-menu-column">
          <QueueList />
        </div>
        <div className="queue-list-column">
          <QueueTable queueType={this.props.match.params.queueType} retrieveMoves={RetrieveMovesForOffice} />
        </div>
      </div>
    );
  }
}

export class RenderWithOrWithoutHeader extends Component {
  render() {
    const Tag = this.props.tag;
    const Component = this.props.component;
    return (
      <>
        {this.props.withHeader && <QueueHeader />}
        <Tag role="main" className="site__content">
          <Component {...this.props} />
        </Tag>
      </>
    );
  }
}

export class OfficeWrapper extends Component {
  state = { hasError: false };

  componentDidMount() {
    document.title = 'Transcom PPP: Office';
    this.props.loadInternalSchema();
    this.props.loadPublicSchema();
    this.props.getCurrentUserInfo();
  }

  componentDidCatch(error, info) {
    this.setState({
      hasError: true,
      error,
      info,
    });
  }

  render() {
    const ConditionalWrap = ({ condition, wrap, children }) => (condition ? wrap(children) : <>{children}</>);
    const { context: { flags: { too } } = { flags: { too: null } } } = this.props;
    const DivOrMainTag = detectIE11() ? 'div' : 'main';
    const { userIsLoggedIn } = this.props;
    return (
      <ConnectedRouter history={history}>
        <div className="Office site">
          {!userIsLoggedIn && <QueueHeader />}
          <ConditionalWrap
            condition={!userIsLoggedIn}
            wrap={children => (
              <DivOrMainTag role="main" className="site__content">
                {children}
              </DivOrMainTag>
            )}
          >
            <LogoutOnInactivity />
            {this.state.hasError && <SomethingWentWrong error={this.state.error} info={this.state.info} />}
            {!this.state.hasError && (
              <Switch>
                <Route
                  exact
                  path="/"
                  component={({ location }) => (
                    <Redirect
                      from="/"
                      to={{
                        ...location,
                        pathname: '/queues/new',
                      }}
                    />
                  )}
                />
                <PrivateRoute
                  path="/queues/:queueType/moves/:moveId"
                  component={props => (
                    <RenderWithOrWithoutHeader component={MoveInfo} withHeader={true} tag={DivOrMainTag} {...props} />
                  )}
                />
                <PrivateRoute
                  path="/queues/:queueType"
                  component={props => (
                    <RenderWithOrWithoutHeader component={Queues} withHeader={true} tag={DivOrMainTag} {...props} />
                  )}
                />
                <PrivateRoute
                  path="/moves/:moveId/orders"
                  component={props => (
                    <RenderWithOrWithoutHeader
                      component={OrdersInfo}
                      withHeader={false}
                      tag={DivOrMainTag}
                      {...props}
                    />
                  )}
                />
                <PrivateRoute
                  path="/moves/:moveId/documents/:moveDocumentId?"
                  component={props => (
                    <RenderWithOrWithoutHeader
                      component={DocumentViewer}
                      withHeader={false}
                      tag={DivOrMainTag}
                      {...props}
                    />
                  )}
                />
                {!isProduction && (
                  <PrivateRoute
                    path="/playground"
                    component={props => (
                      <RenderWithOrWithoutHeader
                        component={ScratchPad}
                        withHeader={true}
                        tag={DivOrMainTag}
                        {...props}
                      />
                    )}
                  />
                )}
                {too && <PrivateRoute path="/too/placeholder" component={TOO} />}
                {too && (
                  <PrivateRoute
                    path="/too/customer/6ac40a00-e762-4f5f-b08d-3ea72a8e4b63/details"
                    component={CustomerDetails}
                  />
                )}
              </Switch>
            )}
          </ConditionalWrap>
        </div>
      </ConnectedRouter>
    );
  }
}

OfficeWrapper.defaultProps = {
  loadInternalSchema: no_op,
  loadPublicSchema: no_op,
};

const mapStateToProps = state => {
  const user = selectCurrentUser(state);
  return {
    swaggerError: state.swaggerInternal.hasErrored,
    userIsLoggedIn: user.isLoggedIn,
  };
};

const mapDispatchToProps = dispatch =>
  bindActionCreators({ loadInternalSchema, loadPublicSchema, getCurrentUserInfo }, dispatch);
export default withContext(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(OfficeWrapper),
);
