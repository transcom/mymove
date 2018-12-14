import { no_op } from 'shared/utils';
import WizardPage from 'shared/WizardPage';
import React, { Component } from 'react';
import Summary from './Summary';
import { connect } from 'react-redux';
import { get } from 'lodash';
import WizardHeader from 'scenes/Moves/WizardHeader';
import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';

class Review extends Component {
  componentDidMount() {
    window.scrollTo(0, 0);
  }
  render() {
    const { pages, pageKey, isHHGPPMComboMove } = this.props;

    return (
      <div>
        {isHHGPPMComboMove && (
          <WizardHeader
            right={
              <ProgressTimeline>
                <ProgressTimelineStep name="Move Setup" completed />
                <ProgressTimelineStep name="Review" current />
              </ProgressTimeline>
            }
          />
        )}
        <WizardPage handleSubmit={no_op} pageList={pages} pageKey={pageKey} pageIsValid={true}>
          <h1>Review</h1>
          <p>You're almost done! Please review your details before we finalize the move.</p>
          <Summary />
        </WizardPage>
      </div>
    );
  }
}

const mapStateToProps = (state, ownProps) => ({
  isHHGPPMComboMove: get(state, 'moves.currentMove.selected_move_type') === 'HHG_PPM',
  ...ownProps,
});

export default connect(mapStateToProps)(Review);
