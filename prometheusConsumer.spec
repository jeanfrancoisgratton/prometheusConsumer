%ifarch aarch64
%global _arch aarch64
%global BuildArchitectures aarch64
%endif

%ifarch x86_64
%global _arch x86_64
%global BuildArchitectures x86_64
%endif

%define debug_package   %{nil}
%define _build_id_links none
%define _name prometheusConsumer
%define _prefix /opt
%define _version 2.01.00
%define _rel 0
#%define _arch x86_64
%define _binaryname prometheusSDSendHost

Name:       prometheusConsumer
Version:    %{_version}
Release:    %{_rel}
Summary:    Prometheus File-based Service Discovery listener and consumer

Group:      monitoring api
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/devops/prometheusFileSDapi

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc
#Requires: sudo
#Obsoletes: vmman1 > 1.140

%description
Prometheus File-based Service Discovery listener and consumer

%prep
%autosetup

%build
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin go build -o %{_sourcedir}/%{_binaryname} .
strip %{_sourcedir}/%{_binaryname}

%clean
rm -rf $RPM_BUILD_ROOT

%pre
exit 0

%install
install -Dpm 0755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post

%preun

%postun

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Fri Sep 27 2024 RPM Builder <builder@famillegratton.net> 2.01.00-0
- Added -printconfig flag, added -ls command (jean-francois@famillegratton.net)
- Moved https post in its own file (jean-francois@famillegratton.net)
- Forgot to bump version here (builder@famillegratton.net)

* Tue Sep 24 2024 RPM Builder <builder@famillegratton.net> 2.00.00-0
- Moved flag parsing around (jean-francois@famillegratton.net)
- Rewrite to send data as JSON (jean-francois@famillegratton.net)

* Fri Sep 20 2024 RPM Builder <builder@famillegratton.net> 1.00.00-0
- new package built with tito

