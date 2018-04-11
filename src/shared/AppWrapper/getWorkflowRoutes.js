import React from 'react';
import { connect } from 'react-redux';

import PrivateRoute from 'shared/User/PrivateRoute';
import WizardPage from 'shared/WizardPage';

import Agreement from 'scenes/Legalese';
import Transition from 'scenes/Moves/Transition';
import MoveType from 'scenes/Moves/MoveTypeWizard';
import PpmSize from 'scenes/Moves/Ppm/PPMSizeWizard';
import PpmWeight from 'scenes/Moves/Ppm/Weight';

const Placeholder = props => {
  return (
    <WizardPage
      handleSubmit={() => undefined}
      pageList={props.pageList}
      pageKey={props.pageKey}
    >
      <h1>Placeholder for {props.title}</h1>
      <h2>{props.description}</h2>
    </WizardPage>
  );
};

const stub = (key, pages) => ({ match }) => (
  <Placeholder pageList={pages} pageKey={key} title={key} />
);

const incompleteServiceMember = props => !props.hasCompleteProfile;
const hasMove = props => props.hasMove;
const hasHHG = props => props.selectedMoveType !== 'PPM';
const hasPPM = props => props.selectedMoveType !== 'HHG';
const isCombo = props => props.selectedMoveType === 'COMBO';
const pages = {
  '/service-member/:id/create': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Create your profile',
  },
  '/service-member/:id/name': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Name',
  },
  '/service-member/:id/contact-info': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Your contact info',
  },
  '/service-member/:id/duty-station': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'current duty station',
  },
  '/service-member/:id/residence-address': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Current residence address',
  },
  '/service-member/:id/backup-mailing-address': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Backup mailing address',
  },
  '/service-member/:id/backup-contacts': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Backup contacts',
  },
  '/service-member/:id/transition': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: "OK, your profile's complete",
  },
  '/orders/:id/': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Tell us about your move orders',
  },
  '/orders/:id/upload': {
    render: stub,
    isInFlow: incompleteServiceMember,
    description: 'Upload your orders',
  },
  '/moves/:moveId': {
    isInFlow: hasMove,
    render: (key, pages) => ({ match }) => (
      <MoveType pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/schedule': {
    render: stub,
    isInFlow: hasHHG,
    description: 'Pick a move date',
  },
  '/moves/:moveId/address': {
    render: stub,
    isInFlow: hasHHG,
    description: 'enter your addresses',
  },

  '/moves/:moveId/ppm-transition': {
    isInFlow: isCombo,
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        <Transition />
      </WizardPage>
    ),
  },
  '/moves/:moveId/ppm-start': {
    isInFlow: state => state.selectedMoveType === 'PPM',
    render: (key, pages) => ({ match }) => (
      <WizardPage handleSubmit={() => undefined} pageList={pages} pageKey={key}>
        <form>
          {' '}
          pickup zip, destination zip, secondary pickup, temp storage?{' '}
        </form>
      </WizardPage>
    ),
  },
  '/moves/:moveId/ppm-size': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmSize pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/ppm-incentive': {
    isInFlow: hasPPM,
    render: (key, pages) => ({ match }) => (
      <PpmWeight pages={pages} pageKey={key} match={match} />
    ),
  },
  '/moves/:moveId/review': {
    render: stub,
    isInFlow: hasMove,
    description: 'Review',
  },
  '/moves/:moveId/agreement': {
    isInFlow: hasMove,
    render: (key, pages) => ({ match }) => {
      return (
        <WizardPage
          handleSubmit={() => undefined}
          pageList={pages}
          pageKey={key}
        >
          <Agreement match={match} />
        </WizardPage>
      );
    },
  },
};
export const getPageList = state =>
  Object.keys(pages).filter(pageKey => {
    const page = pages[pageKey];
    return page.isInFlow(state);
  });

const WorkflowRoutes = props => {
  const pageList = getPageList(props);
  return pageList.map(key => {
    const currPage = pages[key];
    const render = currPage.render(key, pageList, currPage.description);
    return <PrivateRoute exact path={key} key={key} render={render} />;
  });
};

const mapStateToProps = state => ({
  hasCompleteProfile: false,
  selectedMoveType: state.submittedMoves.currentMove
    ? state.submittedMoves.currentMove.selected_move_type
    : null,
  hasMove: Boolean(state.submittedMoves.currentMove),
});
export default connect(mapStateToProps)(WorkflowRoutes);
