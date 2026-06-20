Name:           __APP_NAME__
Version:        __APP_VERSION__
Release:        1%{?dist}
Summary:        IT Asset Inventory Manager

License:        GPL-3.0-or-later
BuildArch:      __RPM_ARCH__

Requires:       webkit2gtk4.1, gtk3

%description
A desktop GUI tool to track computers, phones, tablets,
software licences, Windows keys, antivirus and subscriptions.
Built with Go, Wails and Vue; uses a local SQLite database,
no server required.

%install
install -D -m 0755 %{_sourcedir}/inventory \
    %{buildroot}/usr/local/bin/inventory
install -D -m 0644 %{_sourcedir}/inventory.desktop \
    %{buildroot}/usr/share/applications/inventory.desktop

%files
/usr/local/bin/inventory
/usr/share/applications/inventory.desktop

%changelog
* __DATE__ IT Department <it@example.com> __APP_VERSION__-1
- Automated build via Taskfile
