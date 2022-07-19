import React from 'react';

export default {
  title: 'Global/Typography',
};

// Typography
export const headers = () => (
  <div style={{ padding: '20px' }}>
    <p>h1</p>
    <h1>Public Sans 40/48</h1>
    <p>h2</p>
    <h2>Public Sans 28/34</h2>
    <p>h3</p>
    <h3>Public Sans 22/26</h3>
    <p>h4</p>
    <h4>Public Sans 17/20</h4>
    <p>h5</p>
    <h5>Public Sans 15/21</h5>
    <p>h6</p>
    <h6>Public Sans 13/18</h6>
  </div>
);
export const text = () => (
  <div style={{ padding: '20px' }}>
    <p>p</p>
    <p>
      Public Sans 15/23 Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
      labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
      ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat
      nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id
      est laborum.
    </p>
    <p>p small</p>
    <p>
      <small>
        Public Sans 13/18 Faucibus in ornare quam viverra orci sagittis eu volutpat odio. Felis imperdiet proin
        fermentum leo vel orci. Egestas sed sed risus pretium quam vulputate. Consectetur libero id faucibus nisl. Ipsum
        dolor sit amet consectetur adipiscing elit. Id aliquet lectus proin nibh nisl condimentum id venenatis a.
        Pellentesque pulvinar pellentesque habitant morbi tristique senectus. Mattis vulputate enim nulla aliquet
        porttitor lacus luctus accumsan.
      </small>
    </p>
    <p>
      Bulleted list
      <ul>
        <li>Public Sans 15px / 23px</li>
        <li>USWDS bg-base-darkest</li>
      </ul>
    </p>
  </div>
);
export const links = () => (
  <div style={{ padding: '20px' }}>
    <p>a</p>
    <a href="https://materializecss.com/sass.html">USWDS blue-warm-60v</a>
    <p>a:hover</p>
    <a className="hover" href="https://materializecss.com/sass.html">
      USWDS blue-warm-60v
    </a>
    <p>a:visted</p>
    <a className="visited" href="#">
      USWDS bg-violet-warm-60
    </a>
    <p>a:disabled</p>
    <a href="#" className="disabled">
      This link is disabled
    </a>
    <p>a:focus</p>
    <a href="#" className="focus">
      This link is focused
    </a>
    <p>a small</p>
    <small>
      <a href="https://materializecss.com/sass.html">USWDS blue-warm-60v 14/16</a>
    </small>

    <p>a link light</p>
    <div style={{ background: '#000' }}>
      <a href="#" className="usa-link-light">
        This is link is on a dark background
      </a>
    </div>
  </div>
);
