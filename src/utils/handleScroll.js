const distanceFromTop = window.scrollY;
let newActiveSection;

export default function handleScroll(sections, activeSection, setActiveSection) {
  sections.forEach((section) => {
    const sectionEl = document.querySelector(`#s-${section.id}`);

    if (sectionEl?.offsetTop <= distanceFromTop && sectionEl?.offsetTop + sectionEl?.offsetHeight > distanceFromTop) {
      newActiveSection = section.id;
    }
  });

  if (activeSection !== newActiveSection) {
    setActiveSection(newActiveSection);
  }
}
