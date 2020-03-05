import React from 'react';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';

import { Button, Welcome } from '@storybook/react/demo';

storiesOf('Welcome', module).add('to Storybook', () => <Welcome showApp={linkTo('Button')} />);

storiesOf('Components/Button', module)
  .add('with text', () => <Button onClick={action('clicked')}>Hello Button</Button>)
  .add('with some emoji', () => (
    <Button onClick={action('clicked')}>
      <span role="img" aria-label="so cool">
        😀 😎 👍 💯
      </span>
    </Button>
  ));

  storiesOf('Global Styles/Typography', module)
    .add('Headers', () =>
      <div>
        <p>h1</p><h1>Public Sans 40/48</h1>
        <p>h2</p><h2>Public Sans 28/34</h2>
        <p>h3</p><h3>Public Sans 21/23</h3>
        <p>h4</p><h4>Public Sans 17/23</h4>
        <p>h5</p><h5>Public Sans 15/20</h5>
        <p>h6</p><h6>Public Sans 13/16</h6>
      </div>
    )
    .add('Text', () =>
      <div>
        <p>p</p><p>Public Sans 15/18</p>
        <p>p small</p><small>Public Sans 13/16</small>
      </div>
    )
    .add('Links', () =>
      <div>
        <p>a</p><a href="#">USWDS blue-warm-60v</a>
        <p>a:hover</p><a className="hover" href="#">USWDS blue-warm-60v</a>
        <p>a:visted</p><a className="visited"  href="#">USWDS bg-violet-warm-60</a>
        <p>a:disabled</p><a className="disabled"  href="#">This link is disabled</a>
        <p>a small</p><a href="#">USWDS blue-warm-60v 14/16</a>
      </div>
    );
